## Dataset:
https://www.kaggle.com/datasets/ascanipek/skin-diseases/data
https://www.kaggle.com/datasets/haroonalam16/20-skin-diseases-dataset
https://www.kaggle.com/datasets/pritpal2873/multiple-skin-disease-detection-and-classification?select=Skin+Cancer+Dataset
https://www.kaggle.com/datasets/bilalmanzoor2/fyp-skin-disease-dataset
https://www.kaggle.com/datasets/devdope/synthetic-skin-disease-datasetonly-synthetic?select=Herpes
https://www.kaggle.com/datasets/xuannguyenuet2004/skin-disease-dataset
https://www.kaggle.com/datasets/pacificrm/skindiseasedataset
https://www.kaggle.com/datasets/itachi9604/disease-symptom-description-dataset?select=Symptom-severity.csv

## how to run:
```bash
git clone https://github.com/negative318/grab
cd grab
fastapi run main.py
```
create account in ngrok: https://dashboard.ngrok.com/get-started/setup/windows

open new cmd:

```bash
ngrok config add-authtoken ...
ngrok http http://localhost:8000
```