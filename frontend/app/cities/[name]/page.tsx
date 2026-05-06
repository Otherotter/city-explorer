import CityView from "@/components/CityView";

interface Props {
  params: Promise<{ name: string }>;
}

export default async function CityPage({ params }: Props) {
  const { name } = await params;
  const cityName = name.replace(/-/g, " ");
  return <CityView cityName={cityName} />;
}