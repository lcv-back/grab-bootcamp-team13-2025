import React from 'react';
import ReactMarkdown from 'react-markdown';
import rehypeSanitize from 'rehype-sanitize';
import remarkGfm from 'remark-gfm';
import remarkToc from 'remark-toc';

const mockData = {
  name: "Abacavir allergy",
  markdown: `### Comprehensive Overview of **A53 diffuse large b-cell lymphoma**\n\n---\n\n## 1. Overview & Molecular Features\n- **Subtype:** Part of the larger category of diffuse large B-cell lymphomas.\n- **Genetic feature:** Aneuploidy with TP53 inactivation.\n- **Histology:** Medium-sized to large B cells with a diffuse growth pattern.\n- **Prevalence:** Information not available\n\n---\n\n## 2. Common Signs & Symptoms\n> Not every patient shows all manifestations below; list only those mentioned or clearly implied in the source.\n> If a symptom appears more than once with similar wording, list it only once, choosing the most informative phrasing.\n\nInformation not available\n\n---\n\n## 3. Diagnostics\n\n### For each method list:\n- **Method:** Immunohistochemistry (IHC) analysis\n  - **Main role:** Information not available\n  - **Notes:** Information not available\n\n- **Method:** Morphological assessment\n  - **Main role:** Information not available\n  - **Notes:** Based on a single H & E slide.\n\n- **Method:** Genetic testing\n  - **Main role:** Information not available\n  - **Notes:** Has identified four preferentially mutated genes (TP53, MYD88, SPEN, MYC) in untreated tumor samples.\n\n- **Method:** PET-CT\n  - **Main role:** Diagnosing DLBCL.\n  - **Notes:** Provides high sensitivity and specificity but may not be sufficient for detecting indolent or low-volume disease.\n\n---\n\n## 4. Standard Therapy\n\n### 4.1 First‑line regimen\n- R-CHOP regimen (rituximab, cyclophosphamide, doxorubicin, vincristine, and prednisone). This treatment has a high chance of cure, with approximately 60% of patients being able to be cured with anti-lymphoma therapy.\n\n### 4.2 Salvage / second‑line\n- Autologous stem-cell transplantation (ASCT) is considered a standard of care for chemo-sensitive patients with relapsed DLBCL.\n\n### 4.3 Individualised care\n- Treatment approach can vary depending on individual patient factors and the severity of their condition.\n\n---\n\n## 5. Additional Key Points\n1. More research is needed to determine the most effective diagnostic approach for A53 diffuse large B-cell lymphoma specifically.\n2. Some patients may not respond well to R-CHOP treatment, particularly those with TP53-mutated DLBCL. In such cases, alternative treatments may be necessary.\n3. A team approach may be necessary to ensure the best possible outcome.\n\n---\n\n## 6.Principal References\n\n- **[1]** Vodicka 2022 – DLBCL represents a curable disease with 60–70% chance of cure with current chemoimmunotherapy.\n- **[2]** Wang 2020 – As a widely recognized standard regimen, R-CHOP is able to cure.\n- **[3]** Shi 2024 – Approximately 60% of patients with DLBCL can be cured with anti-...\n- **[4]** Xu 2022 – ASCT is the standard of care for chemo-sensitive patients with relapsed DLBCL.\n\n---\n\n> **Disclaimer:** The information above is for reference only. Diagnosis, prognosis, and treatment must be determined by a qualified medical professional.`
};

export default function Page() {
    return (
      <div className="max-w-4xl mx-auto p-6">
        <h1 className="text-3xl font-bold mb-6">{mockData.name}</h1>
        <div
          className="prose prose-lg max-w-none
            prose-headings:scroll-mt-24
            prose-h2:font-semibold prose-h2:text-2xl prose-h2:mt-8 prose-h2:mb-4
            prose-table:border prose-th:border prose-td:border prose-table:border-collapse"
        >
          <ReactMarkdown
            remarkPlugins={[remarkGfm, remarkToc]}
            rehypePlugins={[rehypeSanitize]}
          >
            {mockData.markdown}
          </ReactMarkdown>
        </div>
      </div>
    );
  }