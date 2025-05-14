'use client';
import Image from 'next/image';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { Mail, Loader2 } from 'lucide-react';

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [message, setMessage] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();
  const BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email) return setError('Please enter your email.');
    setError(null);
    setIsLoading(true);

    try {
      const res = await fetch(`${BASE_URL}/auth/forgot-password`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Unknown error occurred.');

      setMessage('A password reset link has been sent to your email.');
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
        <div className="flex justify-center mb-6 cursor-pointer" onClick={() => router.push('/')}>
          <Image src="/assets/logo.png" alt="iSymptom Logo" width={140} height={40} />
        </div>

        <h2 className="text-xl font-semibold text-[#005a74] mb-4 text-center">Forgot Password</h2>

        <form onSubmit={handleSubmit} className="space-y-4 text-base text-left">
          <label className="block text-gray-600 mb-1">Email</label>
          <div className="flex items-center border border-gray-300 rounded-lg px-3 py-2 focus-within:ring-2 focus-within:ring-[#00BDF9] transition">
            <Mail size={16} className="text-gray-600 mr-2" />
            <input
              name="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com"
              autoComplete="email"
              className="w-full bg-transparent outline-none text-base text-gray-700 placeholder-gray-400"
            />
          </div>

          {error && <p className="text-red-500 text-sm">{error}</p>}
          {message && <p className="text-green-500 text-sm">{message}</p>}

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
                Sending request...
              </>
            ) : (
              'Send Request'
            )}
          </button>
        </form>
      </div>
    </main>
  );
}
