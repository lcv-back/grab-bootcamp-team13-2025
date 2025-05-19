import json
import re
from typing import List, Dict

def parse_symptoms_raw(raw: str) -> List[Dict]:
    """
    Try to parse the whole raw as JSON array; if that fails,
    extract individual JSON objects via regex and parse each.
    """
    raw = raw.strip()
    # 1) Try full JSON load
    try:
        return json.loads(raw)
    except json.JSONDecodeError:
        pass

    # 2) Salvage by extracting {...} blocks
    objs = re.findall(r'\{(?:[^{}]|\n)+?\}', raw)
    symptoms = []
    for o in objs:
        try:
            symptoms.append(json.loads(o))
        except json.JSONDecodeError:
            continue
    return symptoms

# List of your input files
input_files = [
    "0_1400.json",
    "1400_2800.json",
    "2800_4200.json",
    "4200_5000.json",
    "5000_6500.json",
    "6500_8000.json",
    "8000_9500.json",
    "9500_10125.json",
]

merged_data = []

for path in input_files:
    with open(path, "r", encoding="utf-8") as f:
        data = json.load(f)
    for item in data:
        # if there's a broken raw payload, clean & reformat
        if "symptoms_raw" in item:
            raw = item.pop("symptoms_raw")
            # strip markdown fences if present
            raw = re.sub(r"^```json\s*", "", raw)
            raw = re.sub(r"\s*```$", "", raw)
            # parse into a list of symptom objects
            parsed = parse_symptoms_raw(raw)
            item["symptoms"] = parsed
            # mark as ok if we recovered at least one object
            item["status"] = "ok" if parsed else item.get("status", "parse_failed")
        merged_data.append(item)

# write out a single merged, normalized file
with open("merged_test.json", "w", encoding="utf-8") as f:
    json.dump(merged_data, f, ensure_ascii=False, indent=2)

print(f"✅ Merged {len(merged_data)} records into merged.json")
