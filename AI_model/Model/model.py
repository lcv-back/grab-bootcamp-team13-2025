import lancedb
import pyarrow as pa
import torch
from PIL import Image
import os
import json
import numpy as np
import pandas as pd
from sklearn.metrics.pairwise import cosine_similarity
from sentence_transformers import SentenceTransformer
from transformers import CLIPModel, CLIPProcessor
from tqdm import tqdm
import pandas as pd


class Model:
    def __init__(self, link_database):
        self.device = "cuda" if torch.cuda.is_available() else "cpu"
        self.similarity_threshold = 0.99
        self.model = CLIPModel.from_pretrained("openai/clip-vit-large-patch14").to(self.device)
        self.processor = CLIPProcessor.from_pretrained("openai/clip-vit-large-patch14")
        self.text_model = SentenceTransformer("nomic-ai/nomic-embed-text-v1", trust_remote_code=True)

        self.lancedb_instance = lancedb.connect(fr"{link_database}")
        if "image" in self.lancedb_instance.table_names():
            self.database1 = self.lancedb_instance["image"]
        else:
            schema = pa.schema([
                pa.field("image_embedding", pa.list_(pa.float32(), 768)),
                pa.field("name", pa.string()),
                pa.field("path", pa.string())
            ])
            self.lancedb_instance.create_table("image", schema=schema)
            self.database1 = self.lancedb_instance["image"]

        if "symptom" in self.lancedb_instance.table_names():
            self.database2 = self.lancedb_instance["symptom"]
        else:
            schema = pa.schema([
                pa.field("id", pa.string()),
                pa.field("name", pa.string()),
                pa.field("symptom", pa.string()),
                pa.field("symptom_embedding", pa.list_(pa.float32(), 768))
            ])
            self.lancedb_instance.create_table("symptom", schema=schema)
            self.database2 = self.lancedb_instance["symptom"]
        
        df_all_sym = self.database2.to_pandas()[['name','symptom']].drop_duplicates()
        self.symptom_counts = df_all_sym.groupby('name').size().to_dict()


        if "symptom_clip" in self.lancedb_instance.table_names():
            self.database3 = self.lancedb_instance["symptom_clip"]
        else:
            schema = pa.schema([
                pa.field("id", pa.string()),
                pa.field("name", pa.string()),
                pa.field("symptom", pa.string()),
                pa.field("symptom_embedding", pa.list_(pa.float32(), 768))
            ])
            self.lancedb_instance.create_table("symptom_clip", schema=schema)
            self.database3 = self.lancedb_instance["symptom_clip"]


    def add_image_to_database(self, input_data=[], check_existing_path=False,
                            batch_size=500, sim_threshold=0.99):

        if check_existing_path:
            existing_paths = set(self.database1.to_pandas()['path'].tolist())

        batch = []
        batch_embs = []
        batch_paths = []

        for name, path in tqdm(input_data, desc="image process"):
            if check_existing_path and path in existing_paths:
                print(f"Path already exists in database, skipping: {path}")
                continue

            image = Image.open(path)
            inputs = self.processor(images=image, return_tensors="pt").to(self.device)
            with torch.no_grad():
                emb = (self.model
                    .get_image_features(**inputs)
                    .cpu()
                    .squeeze()
                    .numpy()
                    .reshape(1, -1))

            # results = self.database1.search(emb.flatten(), vector_column_name='image_embedding').metric("cosine").limit(1).to_list()

            # if results != []:
            #     similar = 1 - results[0].get('_distance', None)
            #     if similar > 0.99:
            #         if results[0].get('name') == name:
            #             similar_path = results[0].get('path')
            #             print(f"Embedding like {similar_path} ({similar:.3f}), skip: {path}")
            #             continue


            # if batch_embs:
            #     arr = np.vstack(batch_embs) # (n_batch, dim)
            #     sims = cosine_similarity(arr, emb) # (n_batch, 1)
            #     max_sim = sims.max()
            #     if max_sim > sim_threshold:
            #         idx_sim = sims.argmax()
            #         dup_path = batch_paths[idx_sim]
            #         print(f"[Batch] Too similar to {dup_path} (sim={max_sim:.3f}), skipping: {path}")
            #         continue

            batch.append({
                "image_embedding": emb.flatten(),
                "name": name,
                "path": path
            })
            batch_embs.append(emb)
            batch_paths.append(path)
            if check_existing_path:
                existing_paths.add(path)

            if len(batch) >= batch_size:
                self.database1.add(batch)
                print(f"Flushed batch of {len(batch)} items to database.")
                batch.clear()
                batch_embs.clear()
                batch_paths.clear()

        if batch:
            self.database1.add(batch)
            print(f"Flushed final batch of {len(batch)} items to database.")

        print("Done processing all images.")

    def add_symptom_to_database(self, input_data = [], check_existing_symptoms = False, batch_size = 500):
        if check_existing_symptoms:
            existing_symptoms = set(self.database2.to_pandas()['id'].to_list())
        data_to_database = []
        for id, name, symptom_name in tqdm(input_data, "symptom process"):
            
            if check_existing_symptoms:
                if id in existing_symptoms:
                    print(f"symptom already exists in database, skipping: {id}")
                    continue

            symptom_embedding = self.text_model.encode(symptom_name).flatten()

            data_to_database.append({
                "id": id,
                "name": name, 
                "symptom": symptom_name, 
                "symptom_embedding": symptom_embedding
            })

            if check_existing_symptoms:
                existing_symptoms.add(id)
            
            if len(data_to_database) >= batch_size:
                self.database2.add(data_to_database)
                print(f"Flushed batch of {len(data_to_database)} items to database.")
                data_to_database.clear()

        if data_to_database:
            self.database2.add(data_to_database)
            print(f"Flushed final batch of {len(data_to_database)} items to database.")

        print("Done processing all images.")

    def add_symptom_clip_to_database(self, input_data = [], check_existing_symptoms = False, batch_size = 500):
        if check_existing_symptoms:
            existing_symptoms = set(self.database3.to_pandas()['id'].to_list())
        data_to_database = []
        for id, name, symptom_name in tqdm(input_data, "symptom process"):
            
            if check_existing_symptoms:
                if id in existing_symptoms:
                    print(f"symptom already exists in database, skipping: {id}")
                    continue

            inputs = self.processor(text=symptom_name,return_tensors="pt").to(self.device)

            with torch.no_grad():
                emb = self.model.get_text_features(**inputs)    
            symptom_embedding = emb.cpu().squeeze().numpy().astype('float32').tolist()

            data_to_database.append({
                "id": id,
                "name": name, 
                "symptom": symptom_name, 
                "symptom_embedding": symptom_embedding
            })

            if check_existing_symptoms:
                existing_symptoms.add(id)
            
            if len(data_to_database) >= batch_size:
                self.database3.add(data_to_database)
                print(f"Flushed batch of {len(data_to_database)} items to database.")
                data_to_database.clear()

        if data_to_database:
            self.database3.add(data_to_database)
            print(f"Flushed final batch of {len(data_to_database)} items to database.")

        print("Done processing all images.")

    def find_query_by_picture(self, image_path, num_data = 10000, top_k = 5, use_replace = False, metric = "cosine"):

        image = Image.open(image_path)

        inputs = self.processor(images=image, return_tensors="pt").to(self.device)
        with torch.no_grad():
            img_emb = self.model.get_image_features(**inputs)
        img_emb = img_emb.cpu().squeeze().numpy()

        results = (
            self.database1
                .search(img_emb, vector_column_name="image_embedding")
                .metric(metric)
                .limit(num_data)
                .to_pandas()
        )

        df_img = results[["name", "path", "_distance"]]

        if metric.lower() == "cosine":
            df_img["similarity"] = (1 - df_img["_distance"]) * 100
        else:
            df_img["similarity"] = df_img["_distance"] * 100

        df_img["similarity"] = df_img["similarity"].round(2)

        if use_replace:
            clip_hits = self.find_query_by_symptom_clips(image_paths = [image_path], num_data = num_data, top_k = num_data, metric = metric)
            df_clip = pd.DataFrame(clip_hits, columns=["name", "similarity"])
            df_comb = pd.concat([
                df_img[["name", "path", "similarity"]],
                df_clip.assign(path="")[['name','path','similarity']]
            ], ignore_index=True)
            idx = df_comb.groupby("name")["similarity"].idxmax()
            df_best = df_comb.loc[idx, ["name", "path", "similarity"]].reset_index(drop=True)
        else:
            idx = df_img.groupby("name")["similarity"].idxmax()
            df_best = df_img.loc[idx, ["name", "path", "similarity"]].reset_index(drop=True)

      
        df_best = df_best.sort_values(by="similarity", ascending=False).reset_index(drop=True)
        df_best = df_best.head(top_k)
        return list(df_best.itertuples(index=False, name=None))

    def find_query_by_pictures(self, image_paths, num_data = 10000, top_k = 5, use_replace = False, metric = "cosine"):

        if image_paths == []:
            return []

        dfs = []
        for p in image_paths:
            hits = self.find_query_by_picture(image_path = p, num_data = num_data, top_k = num_data, use_replace = use_replace, metric = metric)
            df = pd.DataFrame(hits, columns=["name","path","similarity"])
            dfs.append(df)

        df_all = pd.concat(dfs, ignore_index=True)

        df_sum = (
            df_all
            .groupby("name", as_index=False)["similarity"]
            .sum()
            .rename(columns={"similarity":"sum_similarity"})
        )
        N = len(image_paths)
        df_sum["avg_similarity"] = (df_sum["sum_similarity"] / N).round(2)

        idx = df_sum.groupby("name")["avg_similarity"].idxmax()
        df_best = (
            df_sum
            .loc[idx, ["name" , "avg_similarity"]]
            .sort_values("avg_similarity", ascending=False)
            .reset_index(drop=True)
        )
        df_top = df_best.head(top_k)
        return list(df_top.itertuples(index=False, name=None))

    def find_query_by_symptom(self, symptom, num_data=10000, top_k = 5, metric = "cosine"):
        if self.database2 is None:
            print("Database is empty!")
            return []
        symptom_embedding = self.text_model.encode(symptom).flatten()
        results = self.database2.search(symptom_embedding, vector_column_name="symptom_embedding").metric(metric).limit(num_data).to_pandas()
        df = results[["name", "symptom", "_distance"]].copy()
        if metric.lower() == "cosine":
            df["similarity"] = (1 - df["_distance"]) * 100
        else:
            df["similarity"] = df["_distance"] * 100

        df["similarity"] = df["similarity"].round(2)
        idx = df.groupby("name")["similarity"].idxmax()
        df_best = df.loc[idx, ["name", "symptom", "similarity"]].reset_index(drop=True)

        df_best = df_best.sort_values("similarity", ascending=False).reset_index(drop=True)
        
        df_best = df_best.head(top_k)

        return list(df_best.itertuples(index=False, name=None))

    def find_query_by_symptoms(self, symptoms, num_data=10000, top_k=5, metric="cosine", convert_accuracy = True):
        if not symptoms:
            return []

        dfs = []
        for p in symptoms:
            hits = self.find_query_by_symptom(
                symptom=p, num_data=num_data, top_k=num_data, metric=metric
            )
            df = pd.DataFrame(hits, columns=["name", "symptom", "similarity"])
            dfs.append(df)
        df_all = pd.concat(dfs, ignore_index=True)

        df_sum = (
            df_all
            .groupby("name", as_index=False)["similarity"]
            .sum()
            .rename(columns={"similarity": "sum_similarity"})
        )

        N = len(symptoms)

        df_sum["avg_initial"] = (df_sum["sum_similarity"] / N).round(2)

        df_topk = (
            df_sum
            .sort_values("avg_initial", ascending=False)
            .head(top_k)
            .copy()
        )


        if convert_accuracy:
            df_topk["final_score"] = df_topk.apply(
                lambda row: round(
                    row["avg_initial"]* (1 if N > self.symptom_counts.get(row["name"], N) else N / self.symptom_counts.get(row["name"], N)),2),
                axis=1
            )
        else:
            df_topk["final_score"] = df_topk["avg_initial"]

        df_topk = df_topk.sort_values("final_score", ascending=False).reset_index(drop=True)
        return list(df_topk[["name", "final_score"]].itertuples(index=False, name=None))

    def find_query_by_symptom_clip(self, image, num_data = 10000, top_k = 5, metric = "cosine"):
        if self.database3 is None:
            print("Database is empty!")
            return []
        
        image = Image.open(image).convert("RGB")
        inputs = self.processor(images=image, return_tensors="pt").to(self.device)
        with torch.no_grad():
            img_emb = self.model.get_image_features(**inputs)
        img_emb = img_emb.cpu().squeeze().numpy()

        results = self.database3.search(img_emb, vector_column_name="symptom_embedding").metric(metric).limit(num_data).to_pandas()

        df = results[["name", "symptom", "_distance"]].copy()
        if metric.lower() == "cosine":
            df["similarity"] = (1 - df["_distance"]) * 100
        else:
            df["similarity"] = df["_distance"]
        df["similarity"] = df["similarity"].round(2)

        idx = df.groupby("name")["similarity"].idxmax()
        df_best = df.loc[idx, ["name", "symptom", "similarity"]].reset_index(drop=True)

        df_best = df_best.sort_values("similarity", ascending=False).reset_index(drop=True).head(top_k)

        return list(df_best.itertuples(index=False, name=None))

    def find_query_by_symptom_clips(self, image_paths, num_data = 10000, top_k = 5, metric = "cosine"):

        if self.database3 is None:
            print("Database is empty!")
            return []


        all_hits = []
        for img_path in image_paths:
            hits = self.find_query_by_symptom_clip( image=img_path, num_data=num_data, top_k=num_data, metric=metric)
            
            for name, symptom, sim in hits:
                # all_hits.append({"name": name, "symptom": symptom, "similarity": sim})
                # print(symptom)
                all_hits.append({"name": name, "similarity": sim})

        if not all_hits:
            return []


        df_all = pd.DataFrame(all_hits)


        N = len(image_paths)
        df_agg = (
            df_all
            .groupby("name", as_index=False)["similarity"]
            .sum()
            .rename(columns={"similarity": "sum_similarity"})
        )
        df_agg["avg_similarity"] = (df_agg["sum_similarity"] / N).round(2)

        idx = df_agg.groupby("name")["avg_similarity"].idxmax()
        df_best = df_agg.loc[idx, ["name", "avg_similarity"]]


        df_top = (
            df_best
            .sort_values("avg_similarity", ascending=False)
            .reset_index(drop=True)
            .head(top_k)
        )

        return list(df_top.itertuples(index=False, name=None))

    def find_query_by_picture_symptom(self, list_picture = None, list_symptom = None, num_data = 10000, top_k = 5, use_replace = False, metric = "cosine", convert_accuracy = True):

        if self.database1 is None or self.database2 is None:
            print("Database is empty!")
            return []

        image_results = self.find_query_by_pictures(image_paths = list_picture, num_data = num_data, top_k = num_data, use_replace = True, metric = metric)
        symptom_results = self.find_query_by_symptoms(symptoms = list_symptom, num_data = num_data, top_k = num_data, metric = metric, convert_accuracy=False)

        img_dict  = {name: sim for name, sim in image_results}
        sym_dict  = {name: sim for name, sim in symptom_results}
        all_names = set(img_dict) | set(sym_dict)
        combined = []
        for name in all_names:
            avg_img = img_dict.get(name, 0.0)
            avg_sym = sym_dict.get(name, 0.0)
            combined_score = (avg_img + avg_sym) / 2.0
            combined.append((name, avg_img, avg_sym, combined_score))
        
        combined.sort(key=lambda x: x[3], reverse=True)
        topk = combined[:top_k]
        if convert_accuracy:
            N = len(list_symptom)
            corrected = []
            for name, avg_img, avg_sym, score in topk:
                actual = self.symptom_counts.get(name, N)

                factor = N / actual if actual > N else 1.0
                
                avg_sym = round(avg_sym * factor, 2)
                score = (avg_img + avg_sym) / 2.0
                corrected.append((name, avg_img, avg_sym, score))

            corrected.sort(key=lambda x: x[3], reverse=True)
            return corrected[:top_k]

        return topk