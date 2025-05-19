'use client';

import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { User, Mail, Calendar, Smile, Pencil, Loader2, ArrowLeft } from 'lucide-react';
import Image from 'next/image';
import axios from 'axios';

export default function ProfilePage() {
  const { user, loading } = useAuth();
  const router = useRouter();

  const [editMode, setEditMode] = useState(false);
  const [fullname, setFullname] = useState('');
  const [birthday, setBirthday] = useState('');
  const [gender, setGender] = useState('');
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);

  useEffect(() => {
    if (!loading && !user) {
      router.replace('/login');
    }
    if (user) {
      setFullname(user.fullname || '');
      setBirthday(user.birthday ? new Date(user.birthday).toISOString().substr(0, 10) : '');
      setGender(user.gender || '');
    }
  }, [user, loading, router]);

  if (loading) {
    return (
      <main className="flex items-center justify-center min-h-screen bg-gradient-to-br from-[#e0f7ff] to-white">
        <p className="text-gray-600 text-sm">Loading profile...</p>
      </main>
    );
  }

  if (!user) return null;

  const handleSave = async () => {
    const token = localStorage.getItem('access_token');
    if (!fullname.trim() || !birthday || !gender) {
      setError('Please fill in all required fields.');
      return;
    }

    try {
      setSaving(true);
      setError('');
      setSuccess(false);

      await axios.patch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/update-info`,
        {
          fullname,
          birthday: new Date(birthday).toISOString(),
          gender,
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      setSuccess(true);
      setEditMode(false);
    } catch (err: unknown) {
      if (axios.isAxiosError(err)) {
        setError(err.response?.data?.error || 'Failed to update profile.');
      } else {
        setError('Unexpected error occurred. Please try again.');
      }
    } finally {
      setSaving(false);
    }
  };

  return (
    <main className="flex flex-col items-center justify-center min-h-screen bg-gradient-to-br from-[#e0f7ff] to-white p-4 sm:p-6">
      <div className="bg-white shadow-xl rounded-2xl p-6 sm:p-8 w-full max-w-md space-y-6 animate-slide-up">
        <div className="w-full max-w-md mb-4">
          <button
            onClick={() => router.back()}
            className="flex items-center gap-1 text-sm text-[#00BDF9] hover:underline"
          >
            <ArrowLeft size={16} />
            Back
          </button>
        </div>

        <div className="flex justify-center mb-4">
          <Image
            src="/assets/logo.png"
            alt="iSymptom Logo"
            width={140}
            height={40}
            className="cursor-pointer"
            onClick={() => router.push('/')}
          />
        </div>

        <div className="flex justify-between items-center">
          <h2 className="text-2xl font-bold text-[#005a74]">Account Information</h2>
          <button
            onClick={() => setEditMode(!editMode)}
            className="text-[#00BDF9] hover:underline flex items-center"
          >
            <Pencil size={18} className="mr-1" />
            {editMode ? 'Cancel' : 'Edit'}
          </button>
        </div>

        <div className="space-y-4 text-base text-gray-700">
          {/* Fullname */}
          <div className="flex items-center gap-3">
            <User size={20} className="text-[#00BDF9]" />
            {editMode ? (
              <input
                className="w-full bg-transparent outline-none border border-gray-300 rounded-lg px-3 py-2 text-base text-gray-700 placeholder-gray-400 focus-within:ring-2 focus-within:ring-[#00BDF9] transition"
                value={fullname}
                onChange={(e) => setFullname(e.target.value)}
              />
            ) : (
              <span>{fullname}</span>
            )}
          </div>

          {/* Email (readonly) */}
          <div className="flex items-center gap-3">
            <Mail size={20} className="text-[#00BDF9]" />
            <span>{user.email}</span>
          </div>

          {/* Birthday */}
          <div className="flex items-center gap-3">
            <Calendar size={20} className="text-[#00BDF9]" />
            {editMode ? (
              <input
                type="date"
                className="w-full bg-transparent outline-none border border-gray-300 rounded-lg px-3 py-2 text-base text-gray-700 focus-within:ring-2 focus-within:ring-[#00BDF9] transition"
                value={birthday}
                onChange={(e) => setBirthday(e.target.value)}
              />
            ) : (
              <span>{new Date(birthday).toLocaleDateString()}</span>
            )}
          </div>

          {/* Gender */}
          <div className="flex items-center gap-3">
            <Smile size={20} className="text-[#00BDF9]" />
            {editMode ? (
              <select
                className="w-full bg-transparent outline-none border border-gray-300 rounded-lg px-3 py-2 text-base text-gray-700 focus:ring-2 focus:ring-[#00BDF9] transition"
                value={gender}
                onChange={(e) => setGender(e.target.value)}
              >
                <option value="male">Male</option>
                <option value="female">Female</option>
              </select>
            ) : (
              <span>{gender === 'male' ? 'Male' : gender === 'female' ? 'Female' : 'Other'}</span>
            )}
          </div>
        </div>

        {error && <p className="text-red-500 text-sm">{error}</p>}
        {success && <p className="text-green-500 text-sm">Profile updated successfully!</p>}

        {editMode && (
          <button
            onClick={handleSave}
            className="w-full mt-4 flex items-center justify-center bg-[#00BDF9] hover:bg-[#00acd6] text-white py-2 rounded-lg font-semibold transition-colors disabled:opacity-70 disabled:cursor-not-allowed"
            disabled={saving}
          >
            {saving ? (
              <>
                <Loader2 className="animate-spin h-4 w-4 mr-2" />
                Saving...
              </>
            ) : (
              'Save Changes'
            )}
          </button>
        )}
      </div>
    </main>
  );
}
