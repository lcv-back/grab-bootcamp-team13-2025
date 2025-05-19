import { useState, useEffect } from "react";
import { marked } from "marked";
import axios from "axios";
import { Loader2, Copy, X } from "lucide-react";

import { Hospital, Prediction } from "@/types";

interface DiagnosisResultProps {
  predictions: Prediction[];
  onReset: () => void;
  userSymptoms: string[];
  uploadedUrls?: string[];
}

export default function DiagnosisResult({ predictions, onReset, userSymptoms, uploadedUrls }: DiagnosisResultProps) {
  const [activeIndex, setActiveIndex] = useState<number | null>(null);
  const [loadingIndex, setLoadingIndex] = useState<number | null>(null);
  const [diseaseDescriptions, setDiseaseDescriptions] = useState<Record<number, string>>({});
  const [copiedIndex, setCopiedIndex] = useState<number | null>(null);
  const [readMoreIndex, setReadMoreIndex] = useState<number | null>(null);

  const [showHospitalModal, setShowHospitalModal] = useState(false);
  const [hospitalData, setHospitalData] = useState<Hospital[]>([]);
  const [isLoadingHospitals, setIsLoadingHospitals] = useState(false);

  const [confirmResetModal, setConfirmResetModal] = useState(false);


  useEffect(() => {
    if (activeIndex !== null) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "";
    }
  }, [activeIndex]);

  const getTruncatedMarkdown = (text: string) => {
    const splitIndex = text.indexOf("## 2. Common Signs & Symptoms");
    return splitIndex !== -1 ? text.substring(0, splitIndex) : text;
  };

  const handleOpen = async (index: number, name: string) => {
    setActiveIndex(index);
    if (!diseaseDescriptions[index]) {
      setLoadingIndex(index);
      try {
        const res = await axios.get(`${process.env.NEXT_PUBLIC_API_BASE_URL}/disease?disease_name=${encodeURIComponent(name)}`);
        setDiseaseDescriptions(prev => ({ ...prev, [index]: res.data.description }));
      } catch (err) {
        console.error("Error fetching disease description:", err);
        setDiseaseDescriptions(prev => ({ ...prev, [index]: "**Error loading description.**" }));
      } finally {
        setLoadingIndex(null);
      }
    }
  };

  const handleCopy = (index: number) => {
    const name = predictions[index]?.disease.name;
    navigator.clipboard.writeText(name || "");
    setCopiedIndex(index);
    setTimeout(() => setCopiedIndex(null), 1500);
  };

  const handleClose = () => {
    setActiveIndex(null);
    setReadMoreIndex(null);
  };

  const handleShowNearbyHospitals = () => {
    if (!navigator.geolocation) {
      alert("Geolocation is not supported by your browser.");
      return;
    }

    navigator.geolocation.getCurrentPosition(async (pos) => {
      const { latitude, longitude } = pos.coords;

      try {
        setIsLoadingHospitals(true);
        setShowHospitalModal(true);

        const token = localStorage.getItem("access_token");

        const res = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/hospitals/nearest`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({ latitude, longitude }),
        });

        const data = await res.json();
        setHospitalData(data.hospitals || []);
      } catch (err) {
        console.error("Failed to fetch hospital data:", err);
      } finally {
        setIsLoadingHospitals(false);
      }
    }, (err) => {
      alert("Permission denied or failed to get location.");
      console.error("Geolocation error:", err);
    });
  };

  return (
    <>
      <div className="bg-white rounded-3xl shadow-lg p-6 sm:p-8 w-full max-w-xl space-y-6 mx-auto">
        <h2 className="text-lg md:text-xl font-bold text-[#005a74] mb-2">Top 5 Possible Conditions</h2>
        <ul className="space-y-2">
          {predictions.map((p, idx) => (
            <li
              key={idx}
              className="bg-blue-50 px-4 py-3 rounded-xl cursor-pointer hover:shadow-md transition"
              onClick={() => handleOpen(idx, p.disease.name)}
            >
              <div className="flex justify-between items-center">
                <span className="font-medium text-[#005a74]">{p.disease.name}</span>
                <span className="text-sm text-gray-600">{p.probability.toFixed(2)}%</span>
              </div>
              <div className="w-full bg-gray-100 rounded-full h-2 mt-2">
                <div className="bg-[#00BDF9] h-2 rounded-full" style={{ width: `${p.probability}%` }} />
              </div>
              <span className="text-xs text-gray-500 italic">Tap to view details</span>
            </li>
          ))}
        </ul>

        <div className="flex flex-col gap-2 mt-4">
          <h3 className="text-lg font-semibold text-[#005a74]">Your Reported Symptoms</h3>
          {uploadedUrls && uploadedUrls.length > 0 && (
            <div className="">
              <p className="text-sm font-medium text-[#005a74] mb-2">Uploaded Images:</p>
              <div className="flex gap-4 flex-wrap">
                {uploadedUrls.map((url, idx) => (
                  <img
                    key={idx}
                    src={url}
                    alt={`Uploaded ${idx + 1}`}
                    className="w-24 h-24 object-cover rounded-xl border"
                  />
                ))}
              </div>
            </div>
          )}

          {userSymptoms.length === 0 ? (
            <p className="text-sm text-gray-500 italic">No symptoms provided.</p>
          ) : (
            <ul className="list-disc pl-5 space-y-1 text-sm text-gray-700">
              {userSymptoms.map((s, idx) => (
                <li key={idx}>{s}</li>
              ))}
            </ul>
          )}
          <p className="text-black text-sm font-semibold">Disclaimer: This information is provided for reference purposes only. For medical consultation or check-ups, please contact the hospital directly.</p>
        </div>

        <div className="flex justify-between pt-4">
          <h1
            onClick={handleShowNearbyHospitals}
            className="py-2 hover:cursor-pointer text-[#00BDF9] rounded-full font-semibold text-sm"
          >
            View nearby hospitals
          </h1>
          <button
            onClick={() => setConfirmResetModal(true)}
            className="px-6 py-2 bg-[#00BDF9] hover:bg-[#00acd6] text-white rounded-full font-semibold text-sm"
          >
            Restart
          </button>
        </div>
      </div>

      {/* MODAL: Nearby Hospitals */}
      {showHospitalModal && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 px-4 backdrop-blur-sm transition-opacity duration-300"
          onClick={() => setShowHospitalModal(false)}
        >
          <div
            className="relative bg-white rounded-3xl max-w-3xl w-full shadow-2xl transform transition-all duration-300 ease-out overflow-hidden animate-fadeIn"
            onClick={(e) => e.stopPropagation()}
          >
            {/* HEADER FIXED */}
            <div className="flex justify-between items-start gap-4 mb-4 px-6 pt-6 sticky top-0 bg-white z-10">
              <h3 className="text-2xl font-bold text-[#005a74] leading-snug">
                Nearby Hospitals
              </h3>
              <button
                onClick={() => setShowHospitalModal(false)}
                className="text-white bg-[#00BDF9] hover:bg-[#00acd6] rounded-full p-2 shadow-lg"
              >
                <X size={18} />
              </button>
            </div>

            {/* CONTENT SCROLLABLE */}
            <div className="max-h-[75vh] overflow-y-auto px-6 pb-6 rounded-b-3xl scrollbar-thin scrollbar-thumb-[#00BDF9]/50 scrollbar-track-transparent custom-scrollbar">
              {isLoadingHospitals ? (
                <div className="flex items-center gap-2 text-sm italic text-gray-400">
                  <Loader2 className="w-4 h-4 animate-spin text-[#00BDF9]" />
                  Loading hospitals...
                </div>
              ) : hospitalData.length === 0 ? (
                <p className="text-sm text-gray-500">No hospitals found nearby.</p>
              ) : (
                hospitalData.map((h, idx) => (
                  <div key={idx} className="border border-gray-200 rounded-xl p-4 mb-4">
                    <h4 className="font-semibold text-[#005a74]">{h.name}</h4>
                    <p className="text-sm text-gray-600">{h.address}</p>
                    <p className="text-sm text-gray-500">📍 {h.distance_km?.toFixed(1)} km — ⏰ {h.opening_hours}</p>
                    <div className="flex justify-between items-center mt-2 text-sm">
                      <a href={h.website} target="_blank" className="text-[#00BDF9] hover:underline">Visit website</a>
                      <button
                        onClick={() => navigator.clipboard.writeText(h.phone)}
                        className="text-gray-500 hover:text-[#00BDF9]"
                      >
                        📞 Copy phone
                      </button>
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>
        </div>
      )}



      {/* MODAL: Disease Description */}
      {activeIndex !== null && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 px-4 backdrop-blur-sm transition-opacity duration-300"
          onClick={handleClose}
        >
          <div
            className="relative bg-white rounded-3xl max-w-3xl w-full shadow-2xl transform transition-all duration-300 ease-out overflow-hidden animate-fadeIn"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex justify-between items-start gap-4 mb-4 px-6 pt-6 sticky top-0 bg-white z-10">
              <h3 className="text-2xl font-bold text-[#005a74] leading-snug">
                {predictions[activeIndex].disease.name}
              </h3>
              <div className="flex gap-2 items-center">
                <button
                  onClick={() => handleCopy(activeIndex)}
                  className="text-sm text-gray-400 hover:text-[#00BDF9] flex items-center gap-1"
                >
                  <Copy size={16} /> {copiedIndex === activeIndex ? "Copied" : "Copy name"}
                </button>
                <button
                  onClick={handleClose}
                  className="text-white bg-[#00BDF9] hover:bg-[#00acd6] rounded-full p-2 shadow-lg"
                >
                  <X size={18} />
                </button>
              </div>
            </div>
            <div className="max-h-[75vh] overflow-y-auto px-6 pb-6 rounded-b-3xl scrollbar-thin scrollbar-thumb-[#00BDF9]/50 scrollbar-track-transparent custom-scrollbar">
              {loadingIndex === activeIndex ? (
                <div className="flex items-center gap-2 text-sm italic text-gray-400">
                  <Loader2 className="w-4 h-4 animate-spin text-[#00BDF9]" /> Loading...
                </div>
              ) : (
                <article className="prose prose-sm sm:prose-base lg:prose-lg max-w-none text-[#333] leading-relaxed tracking-normal">
                  <div
                    dangerouslySetInnerHTML={{
                      __html: marked.parse(
                        readMoreIndex === activeIndex
                          ? diseaseDescriptions[activeIndex] || ""
                          : getTruncatedMarkdown(diseaseDescriptions[activeIndex] || "")
                      ),
                    }}
                  />
                  {(diseaseDescriptions[activeIndex] || "").includes("## 2. Common Signs & Symptoms") && (
                    <div className="text-right mt-4">
                      <button
                        onClick={() =>
                          setReadMoreIndex(prev => (prev === activeIndex ? null : activeIndex))
                        }
                        className="text-sm text-gray-500 hover:text-[#00BDF9]"
                      >
                        {readMoreIndex === activeIndex ? "Read less" : "Read more"}
                      </button>
                    </div>
                  )}
                </article>
              )}
            </div>
          </div>
        </div>
      )}

      {confirmResetModal && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 px-4 backdrop-blur-sm"
          onClick={() => setConfirmResetModal(false)}
        >
          <div
            className="relative bg-white rounded-2xl max-w-sm w-full p-6 shadow-xl"
            onClick={(e) => e.stopPropagation()}
          >
            <h3 className="text-lg font-bold text-[#005a74] mb-3">Restart test?</h3>
            <p className="text-sm text-gray-600 mb-5">
              All your answers and results will be lost. Are you sure you want to restart the assessment?
            </p>
            <div className="flex justify-end gap-3">
              <button
                onClick={() => setConfirmResetModal(false)}
                className="px-4 py-1.5 rounded-full border text-sm text-gray-600 hover:bg-gray-100"
              >
                Cancel
              </button>
              <button
                onClick={() => {
                  setConfirmResetModal(false);
                  onReset();
                }}
                className="px-5 py-1.5 rounded-full bg-red-500 hover:bg-red-600 text-white text-sm font-semibold"
              >
                Yes, restart
              </button>
            </div>
          </div>
        </div>
      )}

    </>
  );
}
