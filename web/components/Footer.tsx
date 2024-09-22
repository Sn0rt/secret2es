import React from 'react';
import Image from 'next/image';
import Link from 'next/link';

const Footer: React.FC = () => {
  return (
    <footer className="bg-gray-800 text-white p-2 flex justify-between items-center">
      <p>&copy; {new Date().getFullYear()} sn0rt | BSD-3-Clause License</p>
      <Link
        href="https://external-secrets.io/latest/"
        target="_blank"
        rel="noopener noreferrer"
        className="flex items-center"
      >
        <Image
          src="/eso-round-logo.svg"
          alt="External Secrets Operator"
          width={30}
          height={30}
          className="mr-2"
        />
        <span>External Secrets Operator</span>
      </Link>
    </footer>
  );
};

export default Footer;