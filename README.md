# grab-bootcamp-team13-2025
## how to run:

### Frontend

### Backend

### AI model
```bash
git clone https://github.com/lcv-back/grab-bootcamp-team13-2025.git
cd grab-bootcamp-team13-2025
fastapi run AI_main.py
```
create account in ngrok: https://dashboard.ngrok.com/get-started/setup/windows

open new cmd:

```bash
ngrok config add-authtoken ...
ngrok http http://localhost:8000
```

## Dataset:

### Dataset image

https://www.kaggle.com/datasets/ascanipek/skin-diseases/data
https://www.kaggle.com/datasets/haroonalam16/20-skin-diseases-dataset
https://www.kaggle.com/datasets/pritpal2873/multiple-skin-disease-detection-and-classification?select=Skin+Cancer+Dataset
https://www.kaggle.com/datasets/bilalmanzoor2/fyp-skin-disease-dataset
https://www.kaggle.com/datasets/devdope/synthetic-skin-disease-datasetonly-synthetic?select=Herpes
https://www.kaggle.com/datasets/xuannguyenuet2004/skin-disease-dataset
https://www.kaggle.com/datasets/pacificrm/skindiseasedataset
https://www.kaggle.com/datasets/itachi9604/disease-symptom-description-dataset?select=Symptom-severity.csv


## 2. Setup Backend
open new cmd:
## 2.1 Install dependencies:

```bash
cd backend
go mod download
```
## 2.2. Setup environment:
```bash
cd cp .env.example .env
```

Contents of the `.env` file:

```
JWT_SECRET=
MYSQL_ROOT_PASSWORD=
MYSQL_USER=
MYSQL_PASSWORD=
DB_HOST=
DB_PORT=
DB_NAME=isymptom
DB_USER=
DB_PASSWORD=

EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USERNAME=your_email
EMAIL_PASSWORD=<ask gpt how to create an app password>
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
REDIS_URL=redis:6379

MINIO_ENDPOINT=http://localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=isymptom-bucket
MINIO_USE_SSL=false
GIN_MODE=
SENDGRID_API_KEY=

VPS_USER=root
VPS_HOST=

DOCKERHUB_USERNAME=
DOCKERHUB_TOKEN=
```

## 2.3. Init database:
```bash
mysql -u root -p
```

## 2.4. Run server:
```bash
docker-compose up --build
```
