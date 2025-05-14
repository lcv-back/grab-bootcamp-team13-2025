'use client';
import { Suspense, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Lock, CheckCircle, Loader2 } from "lucide-react";
import Image from "next/image";
import axios from "axios";

function ResetPasswordForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const token = searchParams.get("token");

  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  const BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (password.length < 8) {
      setError("Password must be at least 8 characters long.");
      return;
    }

    if (password !== confirm) {
      setError("Passwords do not match.");
      return;
    }

    if (!token) {
      setError("Invalid token.");
      return;
    }

    try {
      setLoading(true);

      await axios.post(`${BASE_URL}/auth/reset-password`, {
        token,
        new_password: password,
      });

      setSuccess(true);
      setTimeout(() => router.push("/login"), 3000);
    } catch (err: unknown) {
      if (axios.isAxiosError(err)) {
        setError(err.response?.data?.error || "Failed to connect to server.");
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("An unknown error occurred.");
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white shadow-xl rounded-2xl p-6 sm:p-8 w-full max-w-md text-center animate-slide-up">
      <div className="flex justify-center mb-6 cursor-pointer" onClick={() => router.push('/')}>
        <Image src="/assets/logo.png" alt="iSymptom Logo" width={140} height={40} />
      </div>

      <h2 className="text-xl font-bold text-[#005a74] mb-4">Reset Your Password</h2>

      {success ? (
        <div className="flex flex-col items-center text-green-600">
          <CheckCircle size={48} className="mb-2" />
          <p>Password reset successful!</p>
          <p className="text-sm mt-1">You will be redirected to login shortly...</p>
        </div>
      ) : (
        <form onSubmit={handleSubmit} className="space-y-4 text-base text-left">
          <div>
            <label className="block text-gray-600 mb-1">New Password</label>
            <div className="flex items-center border border-gray-300 rounded-lg px-3 py-2 focus-within:ring-2 focus-within:ring-[#00BDF9] transition">
              <Lock size={16} className="text-gray-600 mr-2" />
              <input
                type="password"
                autoComplete="new-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full bg-transparent outline-none text-base text-gray-700 placeholder-gray-400"
                placeholder="Enter new password"
                required
              />
            </div>
          </div>

          <div>
            <label className="block text-gray-600 mb-1">Confirm Password</label>
            <div className="flex items-center border border-gray-300 rounded-lg px-3 py-2 focus-within:ring-2 focus-within:ring-[#00BDF9] transition">
              <Lock size={16} className="text-gray-600 mr-2" />
              <input
                type="password"
                autoComplete="new-password"
                value={confirm}
                onChange={(e) => setConfirm(e.target.value)}
                className="w-full bg-transparent outline-none text-base text-gray-700 placeholder-gray-400"
                placeholder="Re-enter new password"
                required
              />
            </div>
          </div>

          {error && <p className="text-red-500 text-sm">{error}</p>}

          <button
            type="submit"
            disabled={loading}
            className={`w-full flex items-center justify-center bg-[#00BDF9] hover:bg-[#00acd6] text-white py-2 rounded-lg font-semibold transition-colors ${
              loading ? 'opacity-70 cursor-not-allowed' : ''
            }`}
          >
            {loading ? (
              <>
                <Loader2 className="animate-spin h-4 w-4 mr-2" />
                Submitting...
              </>
            ) : (
              "Submit"
            )}
          </button>
        </form>
      )}
    </div>
  );
}

export default function ResetPasswordPage() {
  return (
    <main className="min-h-screen flex items-center justify-center bg-gradient-to-b from-[#e0f7ff] to-white px-4 py-8">
      <Suspense fallback={<p className="text-gray-500">Loading...</p>}>
        <ResetPasswordForm />
      </Suspense>
    </main>
  );
}
