"use client";

import { useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import dynamic from 'next/dynamic';
const Lottie = dynamic(() => import('lottie-react'), { ssr: false });
import animationData from '@/data/logo-landing.json';

export default function LandingPage() {
  const router = useRouter();
  const [isFading, setIsFading] = useState(false);


  useEffect(() => {
    const timer = setTimeout(() => {
      setIsFading(true);
      setTimeout(() => {
        router.push("/home");
      }, 500); // fade duration
    }, 4500); // animation duration
    return () => clearTimeout(timer);
  }, [router]);

  if (!animationData) return null;

  return (
    <main
      className={`flex flex-col bg-[#f9fafa] items-center justify-center h-screen w-screen p-4 transition-opacity duration-500 ${
        isFading ? "opacity-0" : "opacity-100"
      }`}
    >
      <div className="w-[80%] max-w-[500px]"> 
        <Lottie animationData={animationData} loop={false} />
      </div>
    </main>
  );
}
