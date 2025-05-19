'use client';
import Image from 'next/image';
import { useState, useEffect } from 'react';
import { registerUser, loginUser } from '@/lib/api/auth';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { Lock, User, Calendar, Mail, Loader2 } from 'lucide-react';
import Link from 'next/link';

export default function RegisterPage() {
  const [gender, setGender] = useState<string>('');
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<boolean>(false);
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
    const fullname = form.fullname.value.trim();
    const birthday = form.birthday.value;
    const password = form.password.value;
    const confirm = form.confirm.value;

    if (!email || !fullname || !birthday || !password || !confirm) {
      setError('Please fill in all required fields');
      return;
    }
    if (!gender) return setError('Please select a gender');
    if (password.length < 8) return setError('Password must be at least 8 characters');
    if (password !== confirm) return setError('Password confirmation does not match');

    setIsLoading(true);
    try {
      await registerUser({
        email,
        fullname,
        birthday: new Date(birthday).toISOString(),
        gender,
        password,
      });

      const token = await loginUser({ email, password });
      login(token);
      setSuccess(true);
      setError(null);
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("An unexpected error occurred.");
      }
      setSuccess(false);
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

        <h2 className="text-xl font-semibold text-[#005a74] mb-4 text-center">Create Account</h2>

        <form className="space-y-4 text-base" onSubmit={handleSubmit}>
          <div>
            <label className="block text-gray-600 mb-1">Email</label>
            <div className="flex items-center border border-gray-300 rounded-lg px-3 py-2 focus-within:ring-2 focus-within:ring-[#00BDF9] transition">
              <Mail className="text-gray-600 mr-2" size={16} />
              <input name="email" type="email" placeholder="abc@gmail.com" autoComplete="email" className="w-full bg-transparent outline-none text-base text-gray-700 placeholder-gray-400" />
            </div>
          </div>

          <div>
            <label className="block text-gray-600 mb-1">Full Name</label>
            <div className="flex items-center border border-gray-300 rounded-lg px-3 py-2 focus-within:ring-2 focus-within:ring-[#00BDF9] transition">
              <User className="text-gray-600 mr-2" size={16} />
              <input name="fullname" type="text" placeholder="John Doe" autoComplete="name" className="w-full bg-transparent outline-none text-base text-gray-700 placeholder-gray-400" />
            </div>
          </div>

          <div>
            <label className="block text-gray-600 mb-1">Date of Birth</label>
            <div className="flex items-center border border-gray-300 rounded-lg px-3 py-2 focus-within:ring-2 focus-within:ring-[#00BDF9] transition">
              <Calendar className="text-gray-600 mr-2" size={16} />
              <input name="birthday" type="date" className="w-full bg-transparent outline-none text-base text-gray-700" />
            </div>
          </div>

          <div>
            <label className="block text-gray-600 mb-1">Gender</label>
            <div className="flex items-center space-x-3">
              <button
                type="button"
                onClick={() => setGender('male')}
                className={`px-4 py-2 rounded-full border transition ${
                  gender === 'male' ? 'bg-[#00BDF9] text-white border-transparent' : 'text-gray-600 border-gray-300 hover:bg-gray-100'
                }`}
              >
                Male
              </button>
              <button
                type="button"
                onClick={() => setGender('female')}
                className={`px-4 py-2 rounded-full border transition ${
                  gender === 'female' ? 'bg-[#00BDF9] text-white border-transparent' : 'text-gray-600 border-gray-300 hover:bg-gray-100'
                }`}
              >
                Female
              </button>
            </div>
          </div>

          <div>
            <label className="block text-gray-600 mb-1">Password</label>
            <div className="flex items-center border border-gray-300 rounded-lg px-3 py-2 focus-within:ring-2 focus-within:ring-[#00BDF9] transition">
              <Lock className="text-gray-600 mr-2" size={16} />
              <input name="password" type="password" placeholder="Password" autoComplete="new-password" className="w-full bg-transparent outline-none text-base text-gray-700 placeholder-gray-400" />
            </div>
          </div>

          <div>
            <label className="block text-gray-600 mb-1">Confirm Password</label>
            <div className="flex items-center border border-gray-300 rounded-lg px-3 py-2 focus-within:ring-2 focus-within:ring-[#00BDF9] transition">
              <Lock className="text-gray-600 mr-2" size={16} />
              <input name="confirm" type="password" placeholder="Repeat password" className="w-full bg-transparent outline-none text-base text-gray-700 placeholder-gray-400" />
            </div>
          </div>

          {error && <p className="text-red-500 text-sm">{error}</p>}
          {success && <p className="text-green-500 text-sm">Registration successful!</p>}

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
                Creating account...
              </>
            ) : (
              'Register'
            )}
          </button>
        </form>
      </div>
    </main>
  );
}
