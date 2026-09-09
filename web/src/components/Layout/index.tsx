import { useMemo, useState } from "react";
import {
  Layout,
  Menu,
  Dropdown,
  Avatar,
  Button,
  Space,
  App as AntApp,
  Modal,
  Form,
  Input,
} from "antd";
import type { MenuProps } from "antd";
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  UserOutlined,
  LogoutOutlined,
  IdcardOutlined,
  KeyOutlined,
  FileTextOutlined,
  AppstoreOutlined,
  TeamOutlined,
  SafetyCertificateOutlined,
  BarsOutlined,
  SettingOutlined,
  EditOutlined,
  LinkOutlined,
  DashboardOutlined,
  HomeOutlined,
} from "@ant-design/icons";
import { Outlet, useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "../../store/useAuth";
import { useMenu } from "../../store/useMenu";
import { logout as logoutApi, changePassword } from "../../api/auth";
import type { MenuItem } from "../../api/types";

const { Sider, Header, Content } = Layout;

/** 展示端（博客前台）访问地址：优先取环境变量，开发环境兜底 localhost:5174，生产兜底同源根路径 */
const BLOG_URL =
  import.meta.env.VITE_BLOG_URL ||
  (import.meta.env.DEV ? "http://localhost:5174" : "/");

type AntMenuItem = Required<MenuProps>["items"][number];

/** path → 图标 映射（后端菜单树中无图标字段，按路径约定匹配） */
const iconMap: Record<string, React.ReactNode> = {
  "/dashboard": <DashboardOutlined />,
  "/article": <EditOutlined />,
  "/category": <AppstoreOutlined />,
  "/user": <TeamOutlined />,
  "/role": <SafetyCertificateOutlined />,
  "/menu": <BarsOutlined />,
  "/site": <SettingOutlined />,
  "/profile-config": <IdcardOutlined />,
  "/links": <LinkOutlined />,
};

/** 兜底静态菜单（后端未返回菜单树时使用，全量展示） */
const fallbackMenus: MenuItem[] = [
  {
    id: 0,
    parent_id: 0,
    menu_name: "仪表盘",
    menu_type: 2,
    path: "/dashboard",
    sort: 0,
    status: 1,
  },
  {
    id: 1,
    parent_id: 0,
    menu_name: "文章管理",
    menu_type: 2,
    path: "/article",
    sort: 1,
    status: 1,
  },
  {
    id: 2,
    parent_id: 0,
    menu_name: "分类管理",
    menu_type: 2,
    path: "/category",
    sort: 2,
    status: 1,
  },
  {
    id: 3,
    parent_id: 0,
    menu_name: "权限管理",
    menu_type: 1,
    path: "/system",
    sort: 3,
    status: 1,
    children: [
      {
        id: 31,
        parent_id: 3,
        menu_name: "用户管理",
        menu_type: 2,
        path: "/user",
        sort: 1,
        status: 1,
      },
      {
        id: 32,
        parent_id: 3,
        menu_name: "角色管理",
        menu_type: 2,
        path: "/role",
        sort: 2,
        status: 1,
      },
      {
        id: 33,
        parent_id: 3,
        menu_name: "菜单管理",
        menu_type: 2,
        path: "/menu",
        sort: 3,
        status: 1,
      },
    ],
  },
  {
    id: 4,
    parent_id: 0,
    menu_name: "站点设置",
    menu_type: 2,
    path: "/site",
    sort: 4,
    status: 1,
  },
  {
    id: 5,
    parent_id: 0,
    menu_name: "个人资料",
    menu_type: 2,
    path: "/profile-config",
    sort: 5,
    status: 1,
  },
  {
    id: 6,
    parent_id: 0,
    menu_name: "常用网站",
    menu_type: 2,
    path: "/links",
    sort: 6,
    status: 1,
  },
];

/** 将后端菜单树转换为 antd Menu items（跳过按钮类型，目录递归） */
function toMenuItems(nodes: MenuItem[]): AntMenuItem[] {
  return nodes
    .filter((n) => n.status === 1 && n.menu_type !== 3)
    .map((n) => {
      const children = n.children?.filter(
        (c) => c.status === 1 && c.menu_type !== 3,
      );
      if (children && children.length > 0) {
        return {
          key: n.path || `dir-${n.id}`,
          icon: iconMap[n.path ?? ""] ?? <BarsOutlined />,
          label: n.menu_name,
          children: toMenuItems(children),
        };
      }
      return {
        key: n.path || `page-${n.id}`,
        icon: iconMap[n.path ?? ""] ?? <FileTextOutlined />,
        label: n.menu_name,
      };
    });
}

/** 主布局：Sider 权限菜单 + Header 用户信息 + Content 路由出口 */
export default function MainLayout() {
  const navigate = useNavigate();
  const location = useLocation();
  const { message } = AntApp.useApp();
  const userInfo = useAuth((s) => s.userInfo);
  const menus = useAuth((s) => s.menus);
  const doLogout = useAuth((s) => s.logout);
  const collapsed = useMenu((s) => s.collapsed);
  const toggleCollapsed = useMenu((s) => s.toggleCollapsed);
  const [pwdOpen, setPwdOpen] = useState(false);
  const [pwdLoading, setPwdLoading] = useState(false);
  const [pwdForm] = Form.useForm();

  const menuItems = useMemo(() => {
    const source = menus && menus.length > 0 ? menus : fallbackMenus;
    return toMenuItems(source);
  }, [menus]);

  const selectedKeys = useMemo(() => {
    const path = location.pathname;
    // 文章编辑页高亮文章列表菜单
    return [path.startsWith("/article") ? "/article" : path];
  }, [location.pathname]);

  const defaultOpenKeys = useMemo(() => {
    const open: string[] = [];
    const walk = (nodes: MenuItem[]) => {
      nodes.forEach((n) => {
        if (
          n.children?.some(
            (c) => c.path && location.pathname.startsWith(c.path),
          )
        ) {
          open.push(n.path || `dir-${n.id}`);
        }
        if (n.children) walk(n.children);
      });
    };
    walk(menus && menus.length > 0 ? menus : fallbackMenus);
    return open;
  }, [menus, location.pathname]);

  const onMenuClick: MenuProps["onClick"] = ({ key }) => {
    if (key.startsWith("/")) navigate(key);
  };

  const handleLogout = async () => {
    try {
      await logoutApi();
    } catch {
      // 后端登出失败也继续本地清理
    }
    doLogout();
    message.success("已退出登录");
    navigate("/login");
  };

  // 关闭弹窗同时重置表单，避免 destroyOnHidden 不清理外部 form store 导致明文密码残留
  const closePwdModal = () => {
    setPwdOpen(false);
    pwdForm.resetFields();
  };

  const handleChangePwd = async () => {
    let values;
    try {
      values = await pwdForm.validateFields();
    } catch {
      // 校验未通过，字段错误提示由 Form 自动显示
      return;
    }
    setPwdLoading(true);
    try {
      await changePassword({
        old_password: values.old_password,
        new_password: values.new_password,
      });
      // 改密成功后后端已删除 token 强制下线：先提示并关闭/重置弹窗，再清登录态跳转登录页（无需再调 logoutApi）
      message.success("密码修改成功，请重新登录");
      closePwdModal();
      doLogout();
      navigate("/login");
    } catch {
      // 拦截器已统一提示错误（如原密码错误）
    } finally {
      setPwdLoading(false);
    }
  };

  const userMenu: MenuProps = {
    items: [
      {
        key: "profile",
        icon: <IdcardOutlined />,
        label: "个人中心",
        onClick: () => navigate("/profile"),
      },
      {
        key: "changePwd",
        icon: <KeyOutlined />,
        label: "修改密码",
        onClick: () => setPwdOpen(true),
      },
      { type: "divider" },
      {
        key: "logout",
        icon: <LogoutOutlined />,
        label: "退出登录",
        onClick: handleLogout,
      },
    ],
  };

  return (
    // 使用 h-screen 直接基于视口高度，避免 AntApp 等包装层导致的高度链断裂（侧边栏不满高）
    <Layout className="h-screen">
      <Sider
        width={232}
        collapsed={collapsed}
        trigger={null}
        breakpoint="lg"
        className="admin-sider"
      >
        {/* 品牌区 */}
        <div className="flex items-center gap-3 px-5 h-16 border-b border-solid border-[#2b261f]">
          <div
            className="w-9 h-9 shrink-0 flex items-center justify-center rounded-md text-white font-display text-lg"
            style={{ background: "linear-gradient(135deg, #b45309, #7c2d12)" }}
          >
            Z
          </div>
          {!collapsed && (
            <div className="leading-tight">
              <div className="font-display text-[17px] font-bold text-[#f4efe6] tracking-wide">
                Z-blog后台
              </div>
              <div className="text-[10px] text-[#8a7f6f] tracking-[0.2em] uppercase">
                Blog Admin
              </div>
            </div>
          )}
        </div>
        {/* 菜单区独立滚动：菜单项较多时不影响侧边栏满高布局 */}
        <div className="flex-1 min-h-0 overflow-y-auto">
          <Menu
            theme="dark"
            mode="inline"
            items={menuItems}
            selectedKeys={selectedKeys}
            defaultOpenKeys={defaultOpenKeys}
            onClick={onMenuClick}
            className="mt-2"
          />
        </div>
      </Sider>

      <Layout>
        <Header
          className="h-16 px-0 border-b border-solid border-[#e7e0d2]"
          style={{ background: "#fdfbf6" }}
        >
          <div className="flex items-center justify-between w-full h-full px-5 max-w-[1600px] 4k:max-w-[1800px] mx-auto">
            <Space size="middle">
              <Button
                type="text"
                icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
                onClick={toggleCollapsed}
              />
              <span className="text-sm text-[#6d6357] hidden md:inline-block">
                简简单单的博客~
              </span>
            </Space>

            <Space size="middle">
              <Button
                type="text"
                icon={<HomeOutlined />}
                onClick={() =>
                  window.open(BLOG_URL, "_blank", "noopener,noreferrer")
                }
                className="!text-[#6d6357] hover:!text-[#b45309] hover:!bg-[#f4efe6]"
              >
                查看博客
              </Button>

              <Dropdown menu={userMenu} placement="bottomRight">
                <Space className="cursor-pointer select-none" size="small">
                  <Avatar
                    size={32}
                    src={userInfo?.avatar || undefined}
                    icon={<UserOutlined />}
                    style={{ background: "#b45309" }}
                  />
                  <span className="font-medium">
                    {userInfo?.nickname || userInfo?.username || "—"}
                  </span>
                </Space>
              </Dropdown>
            </Space>
          </div>
        </Header>

        <Content className="flex-1 min-h-0 p-6 overflow-auto">
          <Outlet />
        </Content>
      </Layout>

      <Modal
        title="修改密码"
        open={pwdOpen}
        onOk={handleChangePwd}
        onCancel={closePwdModal}
        okText="确认修改"
        cancelText="取消"
        confirmLoading={pwdLoading}
        destroyOnHidden
      >
        <Form form={pwdForm} layout="vertical" className="mt-4">
          <Form.Item
            name="old_password"
            label="原密码"
            rules={[{ required: true, message: "请输入原密码" }]}
          >
            <Input.Password
              placeholder="请输入原密码"
              autoComplete="current-password"
            />
          </Form.Item>
          <Form.Item
            name="new_password"
            label="新密码"
            rules={[
              { required: true, message: "请输入新密码" },
              { min: 6, max: 32, message: "密码长度为 6-32 位" },
            ]}
          >
            <Input.Password
              placeholder="请输入新密码（至少 6 位）"
              autoComplete="new-password"
            />
          </Form.Item>
          <Form.Item
            name="confirm_password"
            label="确认新密码"
            dependencies={["new_password"]}
            rules={[
              { required: true, message: "请再次输入新密码" },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue("new_password") === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error("两次输入的密码不一致"));
                },
              }),
            ]}
          >
            <Input.Password
              placeholder="请再次输入新密码"
              autoComplete="new-password"
            />
          </Form.Item>
        </Form>
      </Modal>
    </Layout>
  );
}
