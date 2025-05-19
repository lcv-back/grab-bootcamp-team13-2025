import json

# Load data
with open('diseases_with_symptom_codes.json', encoding='utf-8') as f:
    data = json.load(f)

# Print diseases without symptoms
print('Diseases without symptoms:')
for d in data:
    if 'symptoms' not in d:
        print(f'global_index: {d["global_index"]}, disease_name: {d["disease_name"]}') 