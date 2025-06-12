import { useState } from 'react';
import { useListPrefectures, useGetPrefecture } from '@/service/generated/client';

export const TopPage = () => {
  // State to track the selected prefecture code
  const [selectedPrefectureCode, setSelectedPrefectureCode] = useState<string | null>(null);

  // Use the Orval-generated hook to fetch prefectures data
  const { data: prefecturesData, isLoading: isPrefecturesLoading, isError: isPrefecturesError } = useListPrefectures();

  // Use the Orval-generated hook to fetch the selected prefecture with municipalities
  const { 
    data: prefectureData, 
    isLoading: isPrefectureLoading, 
    isError: isPrefectureError 
  } = useGetPrefecture(
    selectedPrefectureCode || '', 
    { 
      query: { 
        enabled: !!selectedPrefectureCode 
      } 
    }
  );

  // Handle prefecture selection
  const handlePrefectureClick = (code: string) => {
    setSelectedPrefectureCode(code);
  };

  // Handle back button click
  const handleBackClick = () => {
    setSelectedPrefectureCode(null);
  };

  if (isPrefecturesLoading) {
    return <div>Loading prefectures...</div>;
  }

  if (isPrefecturesError) {
    return <div>Error loading prefectures</div>;
  }

  // If a prefecture is selected, show its details and municipalities
  if (selectedPrefectureCode) {
    if (isPrefectureLoading) {
      return <div>Loading prefecture details...</div>;
    }

    if (isPrefectureError) {
      return <div>Error loading prefecture details</div>;
    }

    return (
      <div>
        <button type="button" onClick={handleBackClick}>← 戻る</button>
        <h1>{prefectureData?.data?.name} の市町村一覧</h1>
        {prefectureData?.data?.municipalities && prefectureData.data.municipalities.length > 0 ? (
          <ul>
            {prefectureData.data.municipalities.map((municipality) => (
              <li key={municipality.id}>
                {municipality.municipality_name_kanji}
              </li>
            ))}
          </ul>
        ) : (
          <p>市町村が見つかりません</p>
        )}
      </div>
    );
  }

  // Otherwise, show the list of prefectures
  return (
    <div>
      <h1>都道府県一覧</h1>
      <ul>
        {prefecturesData?.data?.map((prefecture) => (
          <li key={prefecture.code}>
            <button
              type="button"
              onClick={() => prefecture.code && handlePrefectureClick(prefecture.code)}
              style={{ 
                background: 'none',
                border: 'none',
                padding: 0,
                color: 'blue',
                textDecoration: 'underline',
                cursor: 'pointer',
                font: 'inherit'
              }}
            >
              {prefecture.name} ({prefecture.code})
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
};
