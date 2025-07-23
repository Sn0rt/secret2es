import type { Metadata } from "next";
import "./globals.css";
import Header from "@/components/Header";
import Footer from "@/components/Footer";
import { Toaster } from 'react-hot-toast';
import { TooltipProvider } from "@/components/ui/tooltip";

export const metadata: Metadata = {
  title: "secret2es",
  description: "Convert Kubernetes Secrets to External Secrets",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="flex flex-col min-h-screen">
        <TooltipProvider>
          <Header />
          <main className="flex-grow">
            {children}
          </main>
          <Footer />
          <Toaster />
        </TooltipProvider>
      </body>
    </html>
  );
}
