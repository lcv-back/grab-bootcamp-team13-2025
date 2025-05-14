import { AuthProvider } from '@/contexts/AuthContext';
import "./globals.css";
import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'iSymptom',
  description: 'AI-powered symptom checker',
  icons: {
    icon: '/favicon.ico',
  },
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="vi">
      <body>
        <AuthProvider>{children}</AuthProvider>
      </body>
    </html>
  );
}
