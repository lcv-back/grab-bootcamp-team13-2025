# main.py
from fastapi import FastAPI, UploadFile, File, HTTPException, Request
from fastapi.middleware.cors import CORSMiddleware
from typing import List
from PIL import Image
import io
from AI_model.Model.model import Model
from AI_model.schemas import  ImageItem, SymptomItem, QueryBySymptom, QueryByPicture, CombinedQuery, Predict
import requests
import os
from AI_model.Model.follow_up_ques.main import SymptomInformationGainSelector


def download_images(urls, download_dir, start_index):
    num = start_index
    os.makedirs(download_dir, exist_ok=True)
    local_paths = []
    for url in urls:
        local_path = os.path.join(download_dir, str(num))
        if not os.path.exists(local_path):
            resp = requests.get(url, timeout=10)
            resp.raise_for_status()
            with open(local_path, "wb") as f:
                f.write(resp.content)
        num += 1
        local_paths.append(local_path)
    return local_paths, num

num_download = 0
app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"]
    )


model = Model(link_database="./AI_model/lancedb")
selector = SymptomInformationGainSelector(
        disease_symptoms_path='./AI_model/Model/follow_up_ques/diseases_with_symptom_codes.json',
        symptom_mapping_path='./AI_model/Model/follow_up_ques/symptom_groups_semantic.json',
)

@app.get("/")
async def home():
    return "home"


    # @app.post("/images/add", summary="Add images to database")
    # async def add_images(items: List[ImageItem]):
    #     try:
    #         data = [(it.name, it.path) for it in items]
    #         model.add_image_to_database2(data, check_existing_path=True)
    #         return {"status": "ok", "added": len(data)}
    #     except Exception as e:
    #         raise HTTPException(status_code=500, detail=str(e))

    # @app.post("/symptoms/add", summary="Add symptoms to database")
    # async def add_symptoms(items: List[SymptomItem]):
    #     try:
    #         data = [(it.id, it.name, it.symptom) for it in items]
    #         model.add_symptom_to_database(data, check_existing_symptoms=True)
    #         return {"status": "ok", "added": len(data)}
    #     except Exception as e:
    #         raise HTTPException(status_code=500, detail=str(e))

    # @app.post("/query/symptom", summary="Query by symptom")
    # async def query_symptom(body: QueryBySymptom):
    #     results = model.find_query_by_symptom(body.query, num_data=body.num_data)
    #     return [{"name": n, "score": s} for n, s in results]

    # @app.post("/query/image", summary="Query by image")
    # async def predict_image(request : Request):

    #     form = await request.form()
    #     print(form)
    #     file : UploadFile = form.get("file")
    #     contents = await file.read()
        
    #     image = Image.open(io.BytesIO(contents))

    #     result = model.find_query_by_picture(image)
    #     return {"prediction": result}


@app.post("/predict", summary="Predict diseases from symptoms and images")
async def predict_Symptom(body: Predict):
    global num_download
    download_dir = "./AI_model/downloaded_images"
    if body.image_paths:
        local_paths, downloaded = download_images(body.image_paths, download_dir, num_download)
        num_download = downloaded

    else:
        local_paths = None
    symptoms = None if body.symptoms == [] else body.symptoms
    top_k = 5
    if body.num_data:
        top_k = body.num_data
    if symptoms and local_paths:
        results = model.find_query_by_picture_symptom(list_picture = local_paths, list_symptom = body.symptoms, 
                                                      num_data = 10000, top_k = top_k, use_replace = True, metric = "cosine")
        results_formatted = [
            {"name": entry[0], "probability": entry[3]}
            for entry in results
        ]
        for i in range(5):
            print("image:", results[i][1], " text:" , results[i][2])
    else:
        if symptoms is None:
            results = model.find_query_by_pictures(image_paths = local_paths, num_data = 10000, top_k = top_k, use_replace = True, metric = "cosine")
        elif local_paths is None:
            results = model.find_query_by_symptoms(symptoms = body.symptoms, num_data = 10000, top_k = top_k, metric = "cosine")
        results_formatted = [
            {"name": entry[0], "probability": entry[1]}
            for entry in results
        ]
    

    raw_top = selector.run(
        predicted_records=results_formatted,
        top_n=100
    )
    for symptom in body.answers:
        print("xx", symptom)

    answers_lower = set(ans.lower() for ans in body.answers)
    symptoms_lower = set(sym.lower() for sym in body.symptoms)

    filtered_top = [
        symptom
        for symptom in raw_top
        if symptom.lower() not in answers_lower
        and symptom.lower() not in symptoms_lower
    ]

    final_top = filtered_top[:12]

    return {
        "user_id": body.user_id,
        "predicted_diseases": results_formatted,
        "top_names": final_top
    }