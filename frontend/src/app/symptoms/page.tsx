"use client";
import { useState, useEffect } from "react";
import Header from "@/components/Header";
import ProgressBar from "@/components/ProgressBar";
import Introduction from "@/components/Introduction";
import SymptomForm from "@/components/SymptomForm";
import FollowupQuestions from "@/components/FollowupQuestions";
import DiagnosisResult from "@/components/DiagnosisResult";
import { useAuth } from "@/contexts/AuthContext";
import type { Symptom, Prediction } from "@/types";

export default function SymptomsPage() {
  const [symptoms, setSymptoms] = useState<Symptom[]>([]);
  const [images, setImages] = useState<File[]>([]);
  const [uploadedUrls, setUploadedUrls] = useState<string[]>([]);
  const [token, setToken] = useState("");
  const { user, loading } = useAuth();
  

  const [step, setStep] = useState<'introduction' | 'symptoms' | 'follow-ups' | 'result'>('introduction');
  const [followUpSymptoms, setFollowUpSymptoms] = useState<string[]>([]);
  const [followUpAnswers, setFollowUpAnswers] = useState<Record<string, 'Yes' | 'No'>>({});
  const [predictions, setPredictions] = useState<Prediction[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [followUpCount, setFollowUpCount] = useState(0);

  const [acceptedTerms, setAcceptedTerms] = useState(false);
  const [agreedPrivacy, setAgreedPrivacy] = useState(false);
  const [showTerms, setShowTerms] = useState(false);

  useEffect(() => {
    const savedToken = localStorage.getItem("access_token");
    if (savedToken) setToken(savedToken);
  }, [loading, user]);

  useEffect(() => {
    window.scrollTo(0, 0);
  }, [step]);

  const handleReceiveTopNames = (topNames: string[], predicted: Prediction[]) => {
    if (!Array.isArray(topNames) || topNames.length === 0) {
      alert("No follow-up symptoms received. Please try again later.");
      return;
    }
    setFollowUpSymptoms(topNames);
    setPredictions(predicted);
    setFollowUpCount(1);
    setStep('follow-ups');
  };

  const handleFinalSubmit = async (formattedAnswers: Record<string, 'Yes' | 'No'>) => {
    if (symptoms.length === 0 && !uploadedUrls) {
      alert("Please provide at least one symptom or upload an image before submitting.");
      return;
    }

    try {
      setIsLoading(true);

      const updatedAnswers: Record<string, 'Yes' | 'No'> = {
        ...followUpAnswers,
        ...formattedAnswers
      };

      const yesSymptoms = Object.entries(updatedAnswers)
        .filter(([, ans]) => ans === 'Yes')
        .map(([key]) => key);

      const res = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/symptoms/predict`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          user_id: user?.id || "1",
          symptoms: [...symptoms.map(s => s.name), ...yesSymptoms],
          image_paths: uploadedUrls,
          num_data: 5,
          answers: updatedAnswers
        })
      });

      if (!res.ok) throw new Error(`API error: ${res.status}`);
      const data = await res.json();

      if (!Array.isArray(data.predicted_diseases) || data.predicted_diseases.length === 0) {
        alert("No predictions received. Please try again later.");
        return;
      }

      const converted: Prediction[] = data.predicted_diseases.map((d: { name: string; probability: number }) => ({
        disease: { name: d.name },
        probability: d.probability,
      }));

      const hasHighConfidence = converted.some(d => d.probability > 90);

      if (hasHighConfidence || followUpCount >= 5) {
        setPredictions(converted);
        setStep('result');
      } else {
        const allAskedSymptoms = new Set([
          ...symptoms.map(s => s.name),
          ...Object.keys(updatedAnswers)
        ]);

        const topNames: string[] = data.top_names || [];
        const newFollowUps = topNames.filter((name: string) => !allAskedSymptoms.has(name));

        if (newFollowUps.length === 0) {
          setPredictions(converted);
          setStep('result');
        } else {
          setFollowUpAnswers(updatedAnswers);
          setFollowUpCount(followUpCount + 1);
          setFollowUpSymptoms(newFollowUps);
          setPredictions(converted);
          setStep('follow-ups');
        }
      }
    } catch (err) {
      console.error('Final prediction request failed:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleReset = () => {
    setSymptoms([]);
    setImages([]);
    setUploadedUrls([]);
    setFollowUpAnswers({});
    setFollowUpSymptoms([]);
    setPredictions([]);
    setFollowUpCount(0);
    setStep('introduction');
  };

  return (
    <main className="min-h-screen bg-[#f9fafa]">
      <Header />
      <ProgressBar current={step} />

      <div className="max-w-6xl mx-auto px-4 sm:px-6 md:px-8 py-10">
        {isLoading ? (
          <div className="flex flex-col items-center justify-center py-20 gap-4 text-center">
            <div className="w-6 h-6 border-4 border-[#00BDF9] border-t-transparent rounded-full animate-spin"></div>
            <p className="text-sm text-gray-500">Loading...</p>
          </div>
        ) : (
          <>
            {step === 'introduction' && (
              <Introduction
                onNext={() => setStep('symptoms')}
                acceptedTerms={acceptedTerms}
                agreedPrivacy={agreedPrivacy}
                setAcceptedTerms={setAcceptedTerms}
                setAgreedPrivacy={setAgreedPrivacy}
                showTerms={showTerms}
                setShowTerms={setShowTerms}
              />
            )}

            {step === 'symptoms' && (
              <SymptomForm
                symptoms={symptoms}
                user={user!}
                setSymptoms={setSymptoms}
                images={images}
                setImages={setImages}
                uploadedUrls={uploadedUrls}
                setUploadedUrls={setUploadedUrls}
                token={token}
                onReceiveTopNames={handleReceiveTopNames}
                onBack={() => setStep('introduction')}
              />

            )}

            {step === 'follow-ups' && (
              <FollowupQuestions
                symptoms={followUpSymptoms}
                answers={followUpAnswers}
                setAnswers={(a) => setFollowUpAnswers({ ...followUpAnswers, ...a })}
                onNext={handleFinalSubmit}
                count={followUpCount + 1}
              />
            )}

            {step === 'result' && (
              <DiagnosisResult
                predictions={predictions}
                onReset={handleReset}
                userSymptoms={[
                  ...symptoms.map(s => s.name),
                  ...Object.entries(followUpAnswers)
                    .filter(([, v]) => v === 'Yes')
                    .map(([k]) => k)
                ]}
                uploadedUrls={uploadedUrls}
              />
            )}
          </>
        )}

      </div>
    </main>
  );
}
