# schemas.py
from pydantic import BaseModel # type: ignore
from typing import List, Optional, Dict

class ImageItem(BaseModel):
    name: str
    path: str

class SymptomItem(BaseModel):
    id: str
    name: str
    symptom: str

class QueryBySymptom(BaseModel):
    query: str
    num_data: Optional[int] = 10

class QueryByPicture(BaseModel):
    path: str
    num_data: Optional[int] = 10

class CombinedQuery(BaseModel):
    picture_path: str
    symptom: str
    num_data: Optional[int] = 10

class Predict(BaseModel):
    user_id: int
    symptoms: List[str]
    image_paths: List[str]
    num_data: Optional[int] = 10
    answers: Dict[str, str] = {}