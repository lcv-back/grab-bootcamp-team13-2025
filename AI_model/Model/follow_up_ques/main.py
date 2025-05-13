import json          # ← import ở đây, 1 lần duy nhất
import numpy as np
import pandas as pd
import sys

class SymptomInformationGainSelector:
    def __init__(
        self,
        predicted_path: str = None,
        disease_symptoms_path: str = None,
        symptom_mapping_path: str = None,
        encoding: str = 'utf-8'
    ):
        self.predicted_path    = predicted_path
        self.disease_symptoms_path = disease_symptoms_path
        self.symptom_mapping_path  = symptom_mapping_path
        self.encoding = encoding
        self.load_data()
        sys.stdout.reconfigure(encoding=self.encoding)

    def load_data(self):

        with open(self.disease_symptoms_path, 'r', encoding=self.encoding) as f:
            diseases_symptoms = json.load(f)
        self.disease_symptom_map = {
            d['disease_name']: {s['symptom_code'] for s in d['symptoms']}
            for d in diseases_symptoms
        }

        with open(self.symptom_mapping_path, 'r', encoding=self.encoding) as f:
            mapping_list = json.load(f)
        self.symptom_mappings = {
            item['symptom_code']: item['canonical']
            for item in mapping_list
        }

    def get_symptom_candidates(self, predicted_df):

        candidates = set()
        for name in predicted_df['name']:
            candidates.update(self.disease_symptom_map.get(name, set()))
        symptom_candidates = list(candidates)
        return symptom_candidates

    def build_matrix(self, predicted_df, symptom_candidates):
        n, m = len(predicted_df), len(symptom_candidates)
        matrix = np.zeros((n, m), int)
        for i, name in enumerate(predicted_df['name']):
            codes = self.disease_symptom_map.get(name, set())
            for j, code in enumerate(symptom_candidates):
                if code in codes:
                    matrix[i, j] = 1
        return matrix

    def compute_entropy(self, p: np.ndarray) -> float:
        return -np.sum(p * np.log2(p + 1e-12))

    def compute_information_gain(self, predicted_df, symptom_candidates, matrix, top_n):
        P = predicted_df["probability"].values
        H0 = self.compute_entropy(P)
        igs = []
        for j, code in enumerate(symptom_candidates):
            mask = matrix[:, j].astype(bool)
            p_yes = P[mask].sum()
            p_no = 1 - p_yes

            H_yes = self.compute_entropy(P[mask] / (p_yes + 1e-12)) if p_yes > 0 else 0
            H_no  = self.compute_entropy(P[~mask] / (p_no + 1e-12)) if p_no  > 0 else 0

            H_s = p_yes * H_yes + p_no * H_no
            igs.append((code, H0 - H_s))

        igs.sort(key=lambda x: x[1], reverse=True)
        top_symptoms = igs[:top_n]
        return top_symptoms

    def run(self, predicted_records, top_n) -> list:

        predicted_df = pd.DataFrame(predicted_records)
        predicted_df["probability"] /= predicted_df["probability"].sum()

        symptom_candidates = self.get_symptom_candidates(predicted_df)

        matrix = self.build_matrix(predicted_df, symptom_candidates)

        top_symptoms = self.compute_information_gain(predicted_df, symptom_candidates, matrix, top_n)
        return [
            self.symptom_mappings.get(code, "Unknown")
            for code, _ in top_symptoms
        ]


if __name__ == '__main__':
    
    predicted_list = [
        {
            "name": "Disease by infectious agent",
            "probability": 82.055
        },
        {
            "name": "Exanthem",
            "probability": 79.78999999999999
        },
        {
            "name": "Vasculitis",
            "probability": 77.475
        },
        {
            "name": "Monkeypox",
            "probability": 77.41499999999999
        },
        {
            "name": "Measles",
            "probability": 77.015
        },
        {
            "name": "Chickenpox",
            "probability": 75.785
        },
        {
            "name": "Malign",
            "probability": 75.04
        },
        {
            "name": "Acne",
            "probability": 72.64500000000001
        },
        {
            "name": "Eczema",
            "probability": 68.785
        },
        {
            "name": "Urticaria Hives",
            "probability": 68.56
        },
        {
            "name": "Atopic dermatitis",
            "probability": 67.685
        },
        {
            "name": "Lichen planus",
            "probability": 66.55
        },
        {
            "name": "Warts Molluscum and other Viral Infections",
            "probability": 65.685
        },
        {
            "name": "Psoriasis",
            "probability": 65.11
        },
        {
            "name": "Pellagra",
            "probability": 64.755
        },
        {
            "name": "Tuberous sclerosis",
            "probability": 63.61
        },
        {
            "name": "Poison Ivy",
            "probability": 63.33
        },
        {
            "name": "Tinea Ringworm Candidiases",
            "probability": 63.09
        },
        {
            "name": "Herpes",
            "probability": 62.86
        },
        {
            "name": "Melanoma",
            "probability": 62.85
        },
        {
            "name": "Acrodermatitis enteropathica",
            "probability": 62.435
        },
        {
            "name": "Sebaceous adenocarcinoma",
            "probability": 62.165000000000006
        },
        {
            "name": "Amyloidosis",
            "probability": 62.025
        },
        {
            "name": "Squamous cell carcinoma",
            "probability": 61.91499999999999
        },
        {
            "name": "Pseudoxanthoma elasticum",
            "probability": 61.065
        },
        {
            "name": "Dermatofibroma",
            "probability": 59.765
        },
        {
            "name": "Nevus",
            "probability": 59.53
        },
        {
            "name": "Birt-hogg-dube syndrome",
            "probability": 55.06
        },
        {
            "name": "Benign",
            "probability": 54.53
        },
        {
            "name": "Actinic keratosis",
            "probability": 54.055
        },
        {
            "name": "Addison's disease",
            "probability": 52.83
        },
        {
            "name": "Acanthosis nigricans",
            "probability": 52.81
        },
        {
            "name": "Basal Cell Carcinoma",
            "probability": 50.81
        },
        {
            "name": "Immune system cancer",
            "probability": 50.0
        },
        {
            "name": "Prostate lymphoma",
            "probability": 50.0
        },
        {
            "name": "Multicentric castleman disease",
            "probability": 50.0
        },
        {
            "name": "Gray zone lymphoma",
            "probability": 50.0
        },
        {
            "name": "Acute myeloid leukemia with bcr-abl1",
            "probability": 50.0
        },
        {
            "name": "Cardia cancer",
            "probability": 50.0
        },
        {
            "name": "Mycobacterium avium complex disease",
            "probability": 50.0
        },
        {
            "name": "Gastric lymphoma",
            "probability": 50.0
        },
        {
            "name": "Epithelioid inflammatory myofibroblastic sarcoma",
            "probability": 50.0
        },
        {
            "name": "Mu chain disease",
            "probability": 50.0
        },
        {
            "name": "Anus lymphoma",
            "probability": 50.0
        },
        {
            "name": "Follicular dendritic cell sarcoma",
            "probability": 50.0
        },
        {
            "name": "Childhood lymphoma",
            "probability": 50.0
        },
        {
            "name": "Human immunodeficiency virus infectious disease",
            "probability": 50.0
        },
        {
            "name": "Granulomatosis with polyangiitis",
            "probability": 50.0
        },
        {
            "name": "B-lymphoblastic leukemia/lymphoma with hypodiploidy",
            "probability": 50.0
        },
        {
            "name": "Childhood b-cell acute lymphoblastic leukemia",
            "probability": 50.0
        },
        {
            "name": "Myeloid and lymphoid neoplasms associated with fgfr1 abnormalities",
            "probability": 50.0
        },
        {
            "name": "Erythroleukemia",
            "probability": 50.0
        },
        {
            "name": "Subacute leukemia",
            "probability": 50.0
        },
        {
            "name": "Hematologic cancer",
            "probability": 50.0
        },
        {
            "name": "Lymphoma",
            "probability": 50.0
        },
        {
            "name": "Tuberculosis",
            "probability": 50.0
        },
        {
            "name": "Acute promyelocytic leukemia with pml-rara",
            "probability": 50.0
        },
        {
            "name": "Malignant epithelial mesothelioma",
            "probability": 50.0
        },
        {
            "name": "Methotrexate-associated lymphoproliferation",
            "probability": 50.0
        },
        {
            "name": "Myxofibrosarcoma",
            "probability": 50.0
        },
        {
            "name": "Eosinophilic pneumonia",
            "probability": 50.0
        },
        {
            "name": "Lymph node benign neoplasm",
            "probability": 50.0
        },
        {
            "name": "Myelofibrosis",
            "probability": 50.0
        },
        {
            "name": "Plasmablastic lymphoma",
            "probability": 50.0
        },
        {
            "name": "Tuberculous oophoritis",
            "probability": 50.0
        },
        {
            "name": "B-lymphoblastic leukemia/lymphoma with tcf3-pbx1",
            "probability": 50.0
        },
        {
            "name": "Breast lymphoma",
            "probability": 50.0
        },
        {
            "name": "Nocardiosis",
            "probability": 50.0
        },
        {
            "name": "Tuberculous epididymitis",
            "probability": 50.0
        },
        {
            "name": "Ureteral lymphoma",
            "probability": 50.0
        },
        {
            "name": "A53 diffuse large b-cell lymphoma",
            "probability": 50.0
        },
        {
            "name": "Myeloid and lymphoid neoplasms with eosinophilia and abnormalities of pdgfra, pdgfrb, and fgfr1",
            "probability": 50.0
        },
        {
            "name": "Heavy chain disease",
            "probability": 50.0
        },
        {
            "name": "Collecting duct carcinoma",
            "probability": 50.0
        },
        {
            "name": "Heart sarcoma",
            "probability": 50.0
        },
        {
            "name": "Adult dermatomyositis",
            "probability": 50.0
        },
        {
            "name": "Richter's syndrome",
            "probability": 50.0
        },
        {
            "name": "Organ system cancer",
            "probability": 50.0
        },
        {
            "name": "Acute cholangitis",
            "probability": 50.0
        },
        {
            "name": "Mediastinum seminoma",
            "probability": 50.0
        },
        {
            "name": "Extrapulmonary tuberculosis",
            "probability": 50.0
        },
        {
            "name": "Splenic marginal zone lymphoma",
            "probability": 50.0
        },
        {
            "name": "Angioimmunoblastic t-cell lymphoma",
            "probability": 50.0
        },
        {
            "name": "Acute monocytic leukemia",
            "probability": 50.0
        },
        {
            "name": "Acquired immunodeficiency syndrome",
            "probability": 50.0
        },
        {
            "name": "Pericardial tuberculosis",
            "probability": 50.0
        },
        {
            "name": "Ezb-myc- diffuse large b-cell lymphoma",
            "probability": 50.0
        },
        {
            "name": "Lung giant cell carcinoma",
            "probability": 50.0
        },
        {
            "name": "Myeloid neoplasm",
            "probability": 50.0
        },
        {
            "name": "Colon lymphoma",
            "probability": 50.0
        },
        {
            "name": "Malignant fibrous histiocytoma",
            "probability": 50.0
        },
        {
            "name": "Anaplastic large cell lymphoma",
            "probability": 50.0
        },
        {
            "name": "Intravascular large b-cell lymphoma",
            "probability": 50.0
        },
        {
            "name": "Bn2 diffuse large b-cell lymphoma",
            "probability": 50.0
        },
        {
            "name": "Carcinosarcoma",
            "probability": 50.0
        },
        {
            "name": "Persistent generalized lymphadenopathy",
            "probability": 50.0
        },
        {
            "name": "Multidrug-resistant tuberculosis",
            "probability": 50.0
        },
        {
            "name": "Thymus lymphoma",
            "probability": 50.0
        },
        {
            "name": "Berylliosis",
            "probability": 50.0
        },
        {
            "name": "Chronic lymphocytic leukemia/small lymphocytic lymphoma",
            "probability": 50.0
        },
        {
            "name": "Lymphatic system cancer",
            "probability": 50.0
        },
        {
            "name": "Extranodal marginal zone lymphoma of mucosa-associated lymphoid tissue",
            "probability": 50.0
        },
        {
            "name": "Coccidioidomycosis",
            "probability": 50.0
        },
        {
            "name": "Pulmonary tuberculosis",
            "probability": 50.0
        },
        {
            "name": "Blastic plasmacytoid dendritic cell neoplasm",
            "probability": 50.0
        },
        {
            "name": "Lymphadenitis",
            "probability": 50.0
        },
        {
            "name": "Hairy cell leukemia",
            "probability": 50.0
        },
        {
            "name": "T-cell non-hodgkin lymphoma",
            "probability": 50.0
        },
        {
            "name": "Diffuse large b-cell lymphoma germinal center b-cell type",
            "probability": 50.0
        },
        {
            "name": "St2 diffuse large b-cell lymphoma",
            "probability": 50.0
        },
        {
            "name": "Chronic eosinophilic leukemia",
            "probability": 50.0
        },
        {
            "name": "Sarcomatoid renal cell carcinoma",
            "probability": 50.0
        },
        {
            "name": "Diffuse large b-cell lymphoma",
            "probability": 50.0
        },
        {
            "name": "Grade iii lymphomatoid granulomatosis",
            "probability": 50.0
        },
        {
            "name": "Large b-cell lymphoma",
            "probability": 50.0
        },
        {
            "name": "Composite lymphoma",
            "probability": 50.0
        }
]


    selector = SymptomInformationGainSelector(
        disease_symptoms_path='./diseases_with_symptom_codes.json',
        symptom_mapping_path='./symptom_groups_semantic.json',
    )
    top_names = selector.run(
        predicted_records=predicted_list,
        top_n=12
    )
    print("✅ Top triệu chứng:", top_names)
