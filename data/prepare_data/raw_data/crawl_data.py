import string
import requests
from bs4 import BeautifulSoup
from tqdm import tqdm
import pandas as pd

BASE = 'https://diseasedb.com'
all_diseases = []

for letter in string.ascii_uppercase:
    url = f'{BASE}/diseases/letter/{letter}'
    resp = requests.get(url)
    resp.raise_for_status()
    soup = BeautifulSoup(resp.text, 'html.parser')
    for a in soup.select('a.block.p-4.border'):
        name = a.get_text(strip=True)
        href = a['href']
        all_diseases.append((name, BASE + href))

# PRINT CHỈ DÙNG TIẾNG ASCII
print(f"Tong so benh tim duoc: {len(all_diseases):,}")

def extract_section(soup, title_keyword):
    header = soup.find('h2', string=lambda t: t and title_keyword.lower() in t.lower())
    if not header:
        return ''
    section = header.find_parent('section')
    content_div = section.find('div', class_='prose')
    return content_div.get_text(separator='\n', strip=True) if content_div else ''

results = []
for name, url in tqdm(all_diseases, desc='Processing diseases'):
    try:
        resp = requests.get(url, timeout=10)
        resp.raise_for_status()
        soup = BeautifulSoup(resp.text, 'html.parser')

        description = extract_section(soup, 'Description')
        signs       = extract_section(soup, 'Signs and Symptoms')
        diagnostic  = extract_section(soup, 'Diagnostic Tests')
        treatment   = extract_section(soup, 'Treatment')

        results.append({
            'name': name,
            'description': description,
            'signs_and_symptoms': signs,
            'diagnostic_test': diagnostic,
            'treatment': treatment,
            'url': url
        })
    except Exception as e:
        print(f'Error processing {url}: {e}')

df = pd.DataFrame(results)
output_path = 'diseases_details.csv'
df.to_csv(output_path, index=False)

# CŨNG DÙNG ASCII CHO DÒNG CUỐI
print(f"Da luu ket qua vao {output_path}")
