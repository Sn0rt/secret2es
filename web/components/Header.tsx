import React from 'react';
import Link from 'next/link';
import { Analytics } from "@vercel/analytics/react"

const Header: React.FC = () => {
  return (
    <header className="bg-gray-800 text-white p-3 flex justify-between items-center">
      <div>
        <h1 className="text-2xl font-bold">secret2es</h1>
        <p className="text-base">Convert AVP Secrets to External Secrets</p>
      </div>
      <Link
        href="https://github.com/Sn0rt/secret2es/issues"
        target="_blank"
        rel="noopener noreferrer"
        className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded text-base"
      >
        Feedback
      </Link>
      <Analytics />
    </header>

  );
};

export default Header;