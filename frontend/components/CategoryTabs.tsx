"use client";

import { CATEGORIES, Category } from "@/types";

interface Props {
  active: Category;
  onChange: (category: Category) => void;
}

export default function CategoryTabs({ active, onChange }: Props) {
  return (
    <div className="flex gap-1 overflow-x-auto py-3 scrollbar-hide">
      {CATEGORIES.map((cat) => (
        <button
          key={cat.slug}
          onClick={() => onChange(cat.slug)}
          className={`
            flex items-center gap-2 px-4 py-2 rounded-lg
            whitespace-nowrap text-sm font-medium transition
            ${active === cat.slug
              ? "bg-blue-600 text-white"
              : "text-gray-400 hover:text-white hover:bg-gray-800"
            }
          `}
        >
          <span>{cat.emoji}</span>
          <span>{cat.label}</span>
        </button>
      ))}
    </div>
  );
}