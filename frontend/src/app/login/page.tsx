"use client";
import Image from 'next/image';
import Link from 'next/link';
import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { loginUser } from '@/lib/api/auth';
import { useAuth } from '@/contexts/AuthContext';
import { Mail, Lock, Loader2 } from 'lucide-react';

export default function LoginPage() {
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();
  const { user, login } = useAuth();

  useEffect(() => {
    if (user) router.replace('/symptoms');
  }, [user, router]);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const form = e.currentTarget;
    const email = form.email.value.trim();
    const password = form.password.value;

    if (!email || !password) {
      setError('Please enter your email and password');
      return;
    }

    setIsLoading(true);
    try {
      const token = await loginUser({ email, password });
      login(token);
      setError(null);
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("Unexpected error occurred. Please try again.");
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <main className="flex items-center justify-center min-h-screen bg-gradient-to-br from-[#e0f7ff] to-white px-4 py-8">
      <div className="w-full max-w-md bg-white shadow-xl rounded-2xl p-6 sm:p-8">
        <Link className="flex justify-center mb-6" href={'/'}>
          <Image src="/assets/logo.png" alt="iSymptom Logo" width={140} height={40} />
        </Link>

        <h2 className="text-xl font-semibold text-[#005a74] mb-4 text-center">Sign In</h2>

        <form onSubmit={handleSubmit} className="space-y-4 text-base">
          <div>
            <label className="block text-gray-600 mb-1">Email</label>
            <div className="flex items-center border border-gray-300 rounded-lg px-3 py-2 focus-within:ring-2 focus-within:ring-[#00BDF9] transition">
              <Mail size={16} className="text-gray-600 mr-2" />
              <input
                name="email"
                type="email"
                placeholder="abc@gmail.com"
                autoComplete="email"
                className="w-full bg-transparent outline-none text-base text-gray-700 placeholder-gray-400"
              />
            </div>
          </div>

          <div>
            <label className="block text-gray-600 mb-1">Password</label>
            <div className="flex items-center border border-gray-300 rounded-lg px-3 py-2 focus-within:ring-2 focus-within:ring-[#00BDF9] transition">
              <Lock size={16} className="text-gray-600 mr-2" />
              <input
                name="password"
                type="password"
                placeholder="**************"
                autoComplete="current-password"
                className="w-full bg-transparent outline-none text-base text-gray-700 placeholder-gray-400"
              />
            </div>
          </div>

          <div className="text-right text-sm">
            <Link href="/forgot-password" className="text-gray-500 hover:text-[#00BDF9]">Forgot password?</Link>
          </div>

          {error && <p className="text-red-500 text-sm">{error}</p>}

          <button
            type="submit"
            disabled={isLoading}
            className={`w-full flex items-center justify-center bg-[#00BDF9] hover:bg-[#00acd6] text-white py-2 rounded-lg font-semibold transition-colors ${
              isLoading ? 'opacity-70 cursor-not-allowed' : ''
            }`}
          >
            {isLoading ? (
              <>
                <Loader2 className="animate-spin mr-2 h-4 w-4" />
                Signing in...
              </>
            ) : (
              'Sign In'
            )}
          </button>
        </form>

        <div className="mt-6 text-sm text-gray-600 text-center">
          Don&apos;t have an account?
          <Link href="/register" className="text-[#00BDF9] font-semibold ml-1">Register now</Link>
        </div>
      </div>
    </main>
  );
}
