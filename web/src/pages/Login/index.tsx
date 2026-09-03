import { useState, useEffect } from 'react';
import { Button, Form, Input, Checkbox, App as AntApp } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { login } from '../../api/auth';
import { useAuth } from '../../store/useAuth';

const REMEMBER_KEY = 'remember_username';

export default function Login() {
  const navigate = useNavigate();
  const { message } = AntApp.useApp();
  const setToken = useAuth((s) => s.setToken);
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm();

  // 初始化：读取记住的用户名
  useEffect(() => {
    const saved = localStorage.getItem(REMEMBER_KEY);
    if (saved) {
      form.setFieldsValue({ username: saved, remember: true });
    }
  }, [form]);

  const onFinish = async (values: { username: string; password: string; remember?: boolean }) => {
    setLoading(true);
    try {
      const res = await login({ username: values.username, password: values.password });
      setToken(res.data.token);

      // 记住我：存储或清除用户名
      if (values.remember) {
        localStorage.setItem(REMEMBER_KEY, values.username);
      } else {
        localStorage.removeItem(REMEMBER_KEY);
      }

      message.success('登录成功');
      navigate('/', { replace: true });
    } catch {
      // 错误提示已由响应拦截器统一处理
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex">
      {/* 左侧品牌区：墨色 + 纸纹 */}
      <div
        className="hidden lg:flex flex-col justify-between w-[46%] p-12 text-[#f4efe6] relative overflow-hidden"
        style={{ background: 'linear-gradient(160deg, #191613 0%, #241d14 60%, #33240f 100%)' }}
      >
        <div className="paper-grain absolute inset-0 opacity-60 pointer-events-none" />

        {/* 装饰性圆环 */}
        <div
          className="absolute -top-20 -right-20 w-80 h-80 rounded-full opacity-[0.04] pointer-events-none"
          style={{ border: '2px solid #f4efe6' }}
        />
        <div
          className="absolute -bottom-32 -left-16 w-96 h-96 rounded-full opacity-[0.03] pointer-events-none"
          style={{ border: '1.5px solid #f4efe6' }}
        />

        <div className="relative login-brand-item" style={{ animationDelay: '0.1s' }}>
          <div className="flex items-center gap-3">
            <div
              className="w-11 h-11 flex items-center justify-center rounded-lg font-display text-xl"
              style={{ background: 'linear-gradient(135deg, #b45309, #7c2d12)' }}
            >
              墨
            </div>
            <span className="tracking-[0.35em] text-sm text-[#c9bda6]">墨阁 · BLOG ADMIN</span>
          </div>
        </div>

        <div className="relative">
          <h1
            className="font-display text-[64px] leading-[1.08] m-0 font-black login-brand-item"
            style={{ animationDelay: '0.25s' }}
          >
            落笔成文，
            <br />
            执<span style={{ color: '#d97706' }}>墨</span>有度。
          </h1>
          <p
            className="mt-6 text-[15px] leading-7 text-[#a99c85] max-w-[400px] login-brand-item"
            style={{ animationDelay: '0.4s' }}
          >
            文章、分类、权限与站点配置的统一管理台。
            <br />
            以克制的秩序，承载持续输出的内容。
          </p>
          {/* 装饰琥珀色短线 */}
          <div
            className="login-decor-line mt-6 h-[2px] rounded-full"
            style={{ background: 'linear-gradient(90deg, #b45309, transparent)' }}
          />
        </div>

        <div
          className="relative flex items-center gap-6 text-[11px] tracking-[0.25em] text-[#7a6f5d] uppercase login-brand-item"
          style={{ animationDelay: '0.55s' }}
        >
          <span>RBAC 权限</span>
          <span className="w-1 h-1 rounded-full bg-[#b45309]" />
          <span>ES 全文检索</span>
          <span className="w-1 h-1 rounded-full bg-[#b45309]" />
          <span>OSS 存储</span>
        </div>
      </div>

      {/* 右侧表单区：暖纸底色 */}
      <div className="flex-1 flex items-center justify-center px-6 py-10 lg:p-8 paper-grain relative">
        <div className="w-full max-w-[380px] login-form-enter">
          {/* 小屏 Logo */}
          <div className="lg:hidden mb-8 flex items-center gap-3">
            <div
              className="w-10 h-10 flex items-center justify-center rounded-lg font-display text-lg text-white"
              style={{ background: 'linear-gradient(135deg, #b45309, #7c2d12)' }}
            >
              墨
            </div>
            <span className="font-display text-xl font-bold text-[#1a1815]">墨阁后台</span>
          </div>

          {/* 标题区 */}
          <div className="mb-10">
            <div className="text-xs tracking-[0.3em] text-[#b45309] mb-3 uppercase">Sign in</div>
            <h2 className="font-display text-[34px] m-0 font-bold text-[#1a1815]">管理员登录</h2>
            <p className="mt-3 text-sm text-[#8a7f6f] leading-relaxed">输入账号与密码，进入内容管理台。</p>
          </div>

          {/* 表单 */}
          <Form
            form={form}
            name="login"
            className="login-form"
            size="large"
            onFinish={onFinish}
            autoComplete="new-password"
            layout="vertical"
          >
            <Form.Item
              name="username"
              rules={[{ required: true, message: '请输入用户名' }]}
              style={{ marginBottom: '20px' }}
            >
              <Input
                prefix={<UserOutlined />}
                placeholder="用户名"
              />
            </Form.Item>

            <Form.Item
              name="password"
              rules={[{ required: true, message: '请输入密码' }]}
              style={{ marginBottom: '8px' }}
            >
              <Input.Password
                prefix={<LockOutlined />}
                placeholder="密码"
              />
            </Form.Item>

            {/* 记住我 */}
            <Form.Item name="remember" valuePropName="checked" style={{ marginBottom: '20px' }}>
              <Checkbox>记住用户名</Checkbox>
            </Form.Item>

            <Form.Item className="mb-0">
              <Button
                type="primary"
                htmlType="submit"
                block
                loading={loading}
                className="font-medium"
              >
                登 录
              </Button>
            </Form.Item>
          </Form>

          <div className="mt-12 text-center text-xs tracking-[0.2em] text-[#b3a78f]">
            © 2026 墨阁 · 仅限授权管理员访问
          </div>
        </div>
      </div>
    </div>
  );
}
