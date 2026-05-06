import { Place } from "@/types";

interface Props {
  place: Place;
}

export default function PlaceCard({ place }: Props) {
  return (
    <div className="bg-gray-900 border border-gray-800 rounded-xl
                    p-4 flex flex-col gap-2 hover:border-gray-600
                    transition">

      <div className="flex items-start justify-between gap-2">
        <h3 className="font-semibold text-white leading-tight">
          {place.name}
        </h3>
        {place.subcategory && (
          <span className="text-xs text-gray-400 bg-gray-800
                           px-2 py-1 rounded-full whitespace-nowrap
                           shrink-0">
            {place.subcategory.split(";")[0]}
          </span>
        )}
      </div>

      {place.address && (
        <p className="text-gray-400 text-sm">
          {place.address}
        </p>
      )}

      {place.neighborhood && (
        <p className="text-gray-500 text-xs">
          {place.neighborhood}
        </p>
      )}

      <div className="flex gap-3 mt-1">
        {place.website && (
          <a
            href={place.website}
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-400 hover:text-blue-300
                       text-xs transition"
          >
            Website →
          </a>
        )}
        {place.phone && (
          <a
            href={`tel:${place.phone}`}
            className="text-gray-400 hover:text-white
                       text-xs transition"
          >
            {place.phone}
          </a>
        )}
      </div>

    </div>
  );
}