"use client";
import { useState, useEffect, useRef, useCallback } from "react";
import type { Symptom } from "@/types"; // import từ types.ts

export default function AutocompleteSymptomInput({
  symptoms,
  setSymptoms,
  token,
}: {
  symptoms: Symptom[];
  setSymptoms: (s: Symptom[]) => void;
  token: string;
}) {
  const [query, setQuery] = useState("");
  const [suggestions, setSuggestions] = useState<Symptom[]>([]);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const debounceTimer = useRef<NodeJS.Timeout | null>(null);

  const fetchSuggestions = useCallback(
    async (q: string) => {
      try {
        const res = await fetch(
          `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/symptoms/autocomplete?query=${encodeURIComponent(q)}`,
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );
        const data = await res.json();

        const freeTextSuggestion: Symptom = {
          id: -1,
          name: q,
        };

        setSuggestions([freeTextSuggestion, ...data]);
      } catch (err) {
        console.error("Error fetching suggestions:", err);
      }
    },
    [token]
  );

  const selectSymptom = (symptom: Symptom) => {
    if (!symptoms.some((s) => s.name.toLowerCase() === symptom.name.toLowerCase())) {
      setSymptoms([...symptoms, symptom]);
    }
    setQuery("");
    setSuggestions([]);
    setSelectedIndex(-1);
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "ArrowDown") {
      e.preventDefault();
      setSelectedIndex((prev) => Math.min(prev + 1, suggestions.length - 1));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setSelectedIndex((prev) => Math.max(prev - 1, 0));
    } else if (e.key === "Enter" && selectedIndex >= 0) {
      e.preventDefault();
      selectSymptom(suggestions[selectedIndex]);
    }
  };

  useEffect(() => {
    if (!query) {
      setSuggestions([]);
      return;
    }

    if (debounceTimer.current) clearTimeout(debounceTimer.current);
    debounceTimer.current = setTimeout(() => {
      fetchSuggestions(query);
    }, 300);
  }, [query, fetchSuggestions]);

  const highlightMatch = (text: string, query: string) => {
    const index = text.toLowerCase().indexOf(query.toLowerCase());
    if (index === -1) return text;
    return (
      <>
        {text.substring(0, index)}
        <strong>{text.substring(index, index + query.length)}</strong>
        {text.substring(index + query.length)}
      </>
    );
  };

  return (
    <div className="relative">
      <input
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder="Search, e.g., headache, fever..."
        className="w-full border rounded-full px-4 py-2 outline-none text-base text-gray-600 transition-all duration-200  hover:border-[#00BDF9]
             focus:ring-2 focus:ring-[#00BDF9] focus:border-[#00BDF9] 
             hover:ring-1 hover:ring-[#00BDF9]"
      />
      {suggestions.length > 0 && (
        <div className="absolute mt-1 w-full max-h-60 overflow-auto z-10 bg-white border border-gray-200 rounded-xl shadow-lg">
          {suggestions.map((sug, index) => (
            <div
              key={sug.id === -1 ? `free-text-${sug.name}` : sug.id}
              className={`px-4 py-2 text-base cursor-pointer transition-colors duration-150 ${
                selectedIndex === index
                  ? "bg-[#00BDF9]/10 text-[#005a74]"
                  : "hover:bg-[#00BDF9]/5 hover:text-[#005a74] text-gray-600"
              }`}
              onClick={() => selectSymptom(sug)}
            >
              {highlightMatch(sug.name, query)}
            </div>
          ))}
        </div>
      )}

    </div>
  );
}
