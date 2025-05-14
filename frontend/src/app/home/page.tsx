'use client';
import Image from 'next/image';
import { useRouter } from 'next/navigation';
import dynamic from 'next/dynamic';
import { useAuth } from '@/contexts/AuthContext'; // import Auth context
const Lottie = dynamic(() => import('lottie-react'), { ssr: false });
import animationData from '@/data/doctor-animation.json';

export default function Home() {
  const router = useRouter();
  const { user, loading } = useAuth();

  return (
    <main className="flex flex-col items-center justify-center min-h-screen bg-[#f9fafa] px-6 py-10">
      {/* Logo */}
      <div className="mb-8">
        <Image src="/assets/logo.png" alt="iSymptom Logo" width={180} height={60} />
      </div>

      {/* Main Content */}
      <div className="flex flex-col md:flex-row items-center gap-8 w-full max-w-6xl">
        {/* Animation */}
        <div className="w-full max-w-[280px] md:max-w-md">
          <Lottie animationData={animationData} loop />
        </div>

        {/* Text Content */}
        <div className="text-center md:text-left flex-1 px-2">
          <h2 className="text-2xl md:text-4xl font-bold text-black mb-3">
            {user ? `Hello, ${user.fullname}. How are you today?` : 'Hello, how are you feeling today?'}
          </h2>

          {!user ?
            <p className="text-gray-600 mb-6 text-base md:text-lg">
              Start by choosing an option:
            </p>:
            <p className="text-gray-600 mb-6 text-base md:text-lg">
              Click to start your health journey:
            </p>
          }

          {!loading && (
            <div className="flex justify-center md:justify-start gap-4">
              {user ? (
                <button
                  className="px-6 py-3 bg-[#00BDF9] hover:bg-[#00acd6] transition-all duration-200 shadow-md rounded-full text-white font-semibold text-base md:text-lg"
                  onClick={() => router.push('/symptoms')}
                >
                  Start
                </button>
              ) : (
                <>
                  <button
                    className="px-6 py-3 bg-[#00BDF9] hover:bg-[#00acd6] transition-all duration-200 shadow-md rounded-full text-white font-medium text-base"
                    onClick={() => router.push('/login')}
                  >
                    Login
                  </button>
                  <button
                    className="px-6 py-3 bg-green-400 hover:bg-green-500 transition-all duration-200 shadow-md rounded-full text-white font-medium text-base"
                    onClick={() => router.push('/register')}
                  >
                    Register
                  </button>
                </>
              )}
            </div>
          )}
        </div>
      </div>
    </main>
  );
}
