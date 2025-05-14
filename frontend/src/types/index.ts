export type User = {
    id?: string;
    email: string;
    fullname: string;
    birthday: string;
    gender: string;
    createdAt?: string;
    updatedAt?: string;
  };

export type Symptom = {
    id?: number;
    name: string;
};

export type Disease ={
  name: string;
  description?: string; // Markdown string
}

export type Prediction = {
  disease: Disease;
  probability: number;
}

export type Hospital = {
  name: string;
  distance_km: number;
  address: string;
  phone: string;
  website: string;
  opening_hours: string;
  lat: number;
  lng: number;
}