import { lazy, Suspense, useMemo, type FC } from "react";
import XMarkdown, { type ComponentProps } from "@ant-design/x-markdown";
import "@ant-design/x-markdown/es/XMarkdown/index.css";
import "@ant-design/x-markdown/themes/light.css";
import "@ant-design/x-markdown/themes/dark.css";
import "./MarkdownView.css";
import { slugify, type Heading } from "../utils/markdown";

/** mermaid 图表引擎体积大，按需动态加载，仅在文章含 mermaid 代码块时下载 */
const Mermaid = lazy(() =>
  import("@ant-design/x").then((m) => ({ default: m.Mermaid })),
);

/** 代码高亮器同样拆出主包，随页面空闲预下载，首次代码块渲染时已大概率就绪 */
const CodeHighlighter = lazy(() =>
  import("@ant-design/x").then((m) => ({ default: m.CodeHighlighter })),
);
// 预热：不阻塞首屏渲染，浏览器空闲时预下载高亮器 chunk
import("@ant-design/x");

interface MarkdownViewProps {
  /** Markdown 原文 */
  content: string;
  /** 根容器附加类名 */
  className?: string;
  /** 标题列表（传入则注入 anchor id，用于目录导航） */
  headings?: Heading[];
}

/** 代码块渲染：mermaid 代码块绘制为图表，其余块级代码语法高亮，行内代码保持默认 */
const Code: FC<ComponentProps> = ({ className, children, lang, block }) => {
  if (typeof children !== "string") {
    return <code className={className}>{children}</code>;
  }
  if (block && lang === "mermaid") {
    return (
      <Suspense fallback={<code className={className}>{children}</code>}>
        <Mermaid
          config={{
            themeVariables: {
              fontSize: "18px",
            },
            flowchart: {
              useMaxWidth: false,
              nodeSpacing: 60,
              rankSpacing: 60,
              padding: 20,
            },
            sequence: {
              useMaxWidth: false,
            },
          }}
        >
          {children}
        </Mermaid>
      </Suspense>
    );
  }
  if (block) {
    // lang 可能携带 fence 参数（如 "go title=xx"），仅取语言名
    const pureLang = lang?.split(/[ \t]/)[0] || "";
    return (
      <Suspense fallback={<code className={className}>{children}</code>}>
        <CodeHighlighter lang={pureLang}>{children}</CodeHighlighter>
      </Suspense>
    );
  }
  return <code className={className}>{children}</code>;
};

/**
 * 文章正文渲染组件：XMarkdown 解析 + 纸面排版样式。
 * 样式集中在同目录 MarkdownView.css，便于独立调整。
 */
export default function MarkdownView({
  content,
  className,
  headings,
}: MarkdownViewProps) {
  const config = useMemo(() => {
    if (!headings || headings.length === 0) return undefined;
    return {
      renderer: {
        heading({ text, depth }: { text: string; depth: number }) {
          // text 可能含 HTML 标签（如 <code>），提取纯文本生成 id
          const plainText = text.replace(/<[^>]+>/g, "");
          const id = slugify(plainText);
          return `<h${depth} id="${id}">${text}</h${depth}>`;
        },
      },
    };
  }, [headings]);

  return (
    <div className={`markdown-view${className ? ` ${className}` : ""}`}>
      <XMarkdown
        content={content}
        components={{ code: Code }}
        config={config}
        className="x-markdown-light"
      />
    </div>
  );
}
