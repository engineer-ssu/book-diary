'use client';

import { useEffect, useState } from 'react';

export default function Home() {
  const [data, setData] = useState<{ message: string } | null>(null);

  useEffect(() => {
    // 로컬과 배포 환경 모두 상대 경로인 /api/hello로 호출 가능 (CORS 없음)
    fetch('/api/hello')
      .then((res) => res.json())
      .then((data) => setData(data));
  }, []);

  return (
    <main className="p-8">
      <h1 className="text-2xl font-bold">도서 일기 프로젝트</h1>
      <p className="mt-4">
        Go API 응답: <span className="font-mono text-blue-600">{data ? data.message : '로딩 중...'}</span>
      </p>
    </main>
  );
}