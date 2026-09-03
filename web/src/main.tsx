import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import { ConfigProvider, App as AntApp } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import 'virtual:uno.css';
import './index.css';
import App from './App';

// 墨色纸感主题：暖纸底色 + 墨黑文字 + 琥珀点缀
const themeConfig = {
  token: {
    colorPrimary: '#b45309',
    colorInfo: '#b45309',
    colorLink: '#b45309',
    colorBgLayout: '#f4efe6',
    colorText: '#292521',
    fontFamily:
      "'IBM Plex Sans', -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif",
    borderRadius: 6,
  },
  components: {
    Layout: {
      siderBg: '#191613',
      headerBg: '#fdfbf6',
      bodyBg: '#f4efe6',
    },
    Menu: {
      darkItemBg: '#191613',
      darkSubMenuItemBg: '#12100e',
      darkItemSelectedBg: '#b45309',
      darkItemHoverBg: '#2b261f',
    },
    Table: {
      headerBg: '#f7f2e9',
      headerColor: '#6d6357',
    },
  },
};

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ConfigProvider locale={zhCN} theme={themeConfig}>
      <AntApp>
        <BrowserRouter>
          <App />
        </BrowserRouter>
      </AntApp>
    </ConfigProvider>
  </StrictMode>
);
