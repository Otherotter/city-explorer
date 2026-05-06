import CitySearch from "@/components/CitySearch";

export default function Home() {
  return (
    <main className="min-h-screen bg-gray-950 text-white">
      <div className="max-w-2xl mx-auto px-4 py-24 flex flex-col items-center gap-8">
        <div className="text-center">
          <h1 className="text-4xl font-bold tracking-tight mb-3">
            City Explorer
          </h1>
          <p className="text-gray-400 text-lg">
            If I moved to a new city tomorrow,
            what would my ideal routine look like?
          </p>
        </div>
        <CitySearch />
      </div>
    </main>
  );
}