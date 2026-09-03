import type { ReactNode } from 'react';

interface PageHeaderProps {
  title: string;
  subtitle?: string;
  extra?: ReactNode;
}

/** 页面标题区：衬线大标题 + 说明 + 右侧操作区 */
export default function PageHeader({ title, subtitle, extra }: PageHeaderProps) {
  return (
    <div className="flex items-end justify-between mb-5 gap-4 flex-wrap">
      <div>
        <div className="w-8 h-[3px] bg-[#b45309] mb-3" />
        <h2 className="font-display text-[26px] m-0 font-bold text-[#1a1815] leading-none">
          {title}
        </h2>
        {subtitle && <p className="mt-2 mb-0 text-sm text-[#8a7f6f]">{subtitle}</p>}
      </div>
      {extra && <div className="flex items-center gap-2">{extra}</div>}
    </div>
  );
}
