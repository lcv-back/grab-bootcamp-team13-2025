// components/ProgressBar.tsx
import { Check } from 'lucide-react';

export default function ProgressBar({ current }: { current: 'introduction' | 'symptoms' | 'follow-ups' | 'result' }) {
  const steps = ['introduction', 'symptoms', 'follow-ups', 'result'] as const;
  const currentIndex = steps.indexOf(current);

  return (
    <div className="w-full bg-[#e8f6fb] py-6 px-4 border-b border-[#cce7f2]">
      <div className="relative flex items-center justify-between max-w-5xl mx-auto">
        {/* connecting line */}
        <div className="absolute top-5 left-0 right-0 h-0.5 bg-[#a8ddf1] z-0" />

        {steps.map((step, idx) => {
          const isActive = idx === currentIndex;
          const isCompleted = idx < currentIndex;

          return (
            <div key={step} className="flex flex-col items-center text-center relative z-10 flex-1">
              <div
                className={`w-11 h-11 rounded-full border-2 flex items-center justify-center text-sm font-semibold mb-2 transition-all duration-200
                  ${isActive ? 'bg-white text-[#00BDF9] border-[#00BDF9] scale-110' :
                    isCompleted ? 'bg-white text-green-500 border-green-500' :
                      'bg-white text-gray-400 border-gray-300'}`}
              >
                {isCompleted ? <Check size={18} /> : idx + 1}
              </div>
              <div
                className={`text-sm capitalize font-medium transition-colors duration-200
                  ${isActive ? 'text-[#00BDF9]' : 'text-gray-500 group-hover:text-[#00BDF9]'}`}
              >
                {step}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}