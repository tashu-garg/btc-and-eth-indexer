import type { LucideProps } from 'lucide-react';

export const BTCIcon = ({ size = 24, ...props }: LucideProps) => (
  <svg
    width={size}
    height={size}
    {...props}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <path d="M6 12h9a3 3 0 0 1 0 6H6" />
    <path d="M6 18V6" />
    <path d="M6 12h8a3 3 0 0 0 0-6H6" />
    <path d="M9 6V3" />
    <path d="M13 6V3" />
    <path d="M9 21v-3" />
    <path d="M13 21v-3" />
  </svg>
);

export const ETHIcon = ({ size = 24, ...props }: LucideProps) => (
  <svg
    width={size}
    height={size}
    {...props}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <path d="m12 1 7 11-7 11-7-11Z" />
    <path d="m12 13 7-4-7-7-7 7 7 4Z" />
    <path d="m12 13v10" />
    <path d="m12 1 7 11-7 4-7-4 7-11Z" />
  </svg>
);
