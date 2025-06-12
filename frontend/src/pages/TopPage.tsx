import { useListPrefectures } from '@/service/generated/client';

export const TopPage = () => {
  // Use the Orval-generated hook to fetch prefectures data
  const { data, isLoading, isError } = useListPrefectures();

  if (isLoading) {
    return <div>Loading prefectures...</div>;
  }

  if (isError) {
    return <div>Error loading prefectures</div>;
  }

  return (
    <div>
      <h1>都道府県一覧</h1>
      <ul>
        {data?.data?.map((prefecture) => (
          <li key={prefecture.code}>
            {prefecture.name} ({prefecture.code})
          </li>
        ))}
      </ul>
    </div>
  );
};
