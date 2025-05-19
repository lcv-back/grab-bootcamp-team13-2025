import json
import re

# Input and output paths
input_path = '0_1400.json'
output_path = 'normalized_0_1400.json'

# Read original JSON
with open(input_path, 'r', encoding='utf-8') as f:
    data = json.load(f)

normalized = []

for entry in data:
    disease = entry.get('disease_name', '<unknown>')
    raw = entry.get('symptoms_raw', '')
    
    # Try to extract JSON array of symptoms
    # Match content between the first '[' and the last ']' inclusive
    start = raw.find('[')
    end = raw.rfind(']')
    if start != -1 and end != -1 and end > start:
        symptoms_json_str = raw[start:end+1]
    else:
        print(f"Warning: no JSON array found for {disease}, skipping.")
        continue

    # Attempt parsing
    try:
        symptoms = json.loads(symptoms_json_str)
    except json.JSONDecodeError as e:
        print(f"Error parsing JSON for {disease}: {e}")
        continue

    normalized.append({
        'disease_name': disease,
        'symptoms': symptoms
    })

# Write normalized output
with open(output_path, 'w', encoding='utf-8') as f:
    json.dump(normalized, f, ensure_ascii=False, indent=2)

print(f"Normalized file created at: {output_path}")
