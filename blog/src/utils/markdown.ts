/** 将标题文本转为 URL 友好的 id（用于锚点定位） */
export function slugify(text: string): string {
  return text
    .toLowerCase()
    .trim()
    .replace(/[\s]+/g, '-')
    .replace(/[^\w\u4e00-\u9fa5-]/g, '') // 保留字母数字中文和连字符
    .replace(/--+/g, '-')
    .replace(/^-+|-+$/g, '');
}

export interface Heading {
  level: number;
  text: string;
  id: string;
}

/** 从 Markdown 原文提取标题列表（h1-h4），用于生成目录导航 */
export function extractHeadings(md: string): Heading[] {
  const headings: Heading[] = [];
  const regex = /^(#{1,4})\s+(.+)$/gm;
  let match: RegExpExecArray | null;
  while ((match = regex.exec(md)) !== null) {
    const level = match[1].length;
    const text = match[2].trim();
    headings.push({ level, text, id: slugify(text) });
  }
  return headings;
}
