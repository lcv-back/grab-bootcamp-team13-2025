'use client';
import Image from 'next/image';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

export default function Header() {
  const { user, logout } = useAuth();
  const router = useRouter();

  return (
    <header className="bg-white shadow-md px-6 py-3 sticky top-0 z-50">
      <div className="flex items-center justify-between max-w-7xl mx-auto">
        {/* Logo */}
        <div className="flex items-center space-x-2 cursor-pointer" onClick={() => router.push('/symptoms')}>
          <Image src="/assets/logo.png" alt="iSymptom Logo" width={140} height={36} />
        </div>

        {/* Right section */}
        {user ? (
          <div className="flex items-center space-x-4">
            <Link href="/profile" className="text-[#005a74] font-medium hover:underline truncate max-w-[140px]">
              {user.fullname}
            </Link>
            <button
              onClick={logout}
              className="bg-red-500 hover:bg-red-600 text-white px-4 py-1.5 rounded-full text-sm font-semibold transition"
            >
              Log Out
            </button>
          </div>
        ) : (
          <button
            onClick={() => router.push('/login')}
            className="bg-[#00BDF9] hover:bg-[#00acd6] text-white px-4 py-1.5 rounded-full text-sm font-semibold transition"
          >
            Sign In
          </button>
        )}
      </div>
    </header>
  );
}
