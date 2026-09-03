import { useEffect, useState } from 'react';
import { Avatar, Button, Card, Form, Input, Spin, Upload, App as AntApp } from 'antd';
import { CameraOutlined, SaveOutlined, UserOutlined } from '@ant-design/icons';
import type { UploadFile } from 'antd';
import PageHeader from '../../components/PageHeader';
import { getProfile, updateProfile } from '../../api/profile';
import { uploadToOss } from '../../api/upload';
import { useAuth } from '../../store/useAuth';
import type { Profile } from '../../api/types';

/** 个人中心：查看与编辑当前登录用户的头像、昵称、简介 */
export default function ProfilePage() {
  const { message } = AntApp.useApp();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [profile, setProfile] = useState<Profile | null>(null);
  const [avatar, setAvatar] = useState('');
  const patchUser = useAuth((s) => s.patchUser);

  useEffect(() => {
    setLoading(true);
    getProfile()
      .then((res) => {
        const data = res.data;
        setProfile(data);
        setAvatar(data?.avatar ?? '');
        form.setFieldsValue({
          nickname: data?.nickname,
          bio: data?.bio,
        });
      })
      .catch(() => undefined)
      .finally(() => setLoading(false));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  /** 头像选择后直接上传 OSS，成功后仅更新预览，随表单一起保存 */
  const handleUpload = async (file: File) => {
    if (!file.type.startsWith('image/')) {
      message.error('请选择图片文件');
      return false;
    }
    setUploading(true);
    try {
      const res = await uploadToOss(file);
      setAvatar(res.data.url);
      message.success('头像上传成功，点击保存后生效');
    } catch {
      /* 拦截器已提示 */
    } finally {
      setUploading(false);
    }
    return false; // 阻止 antd 默认上传
  };

  const handleSave = async () => {
    const values = await form.validateFields();
    setSaving(true);
    try {
      const res = await updateProfile({
        nickname: values.nickname ?? '',
        avatar,
        bio: values.bio ?? '',
      });
      setProfile(res.data);
      patchUser({ nickname: res.data.nickname, avatar: res.data.avatar });
      message.success('个人信息已保存');
    } catch {
      /* 拦截器已提示 */
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="page-container">
      <PageHeader title="个人中心" subtitle="管理你的头像、昵称与个人简介" />

      <Spin spinning={loading}>
        <div className="flex flex-col lg:flex-row gap-5 items-start">
          {/* 左侧：资料预览卡 */}
          <Card variant="borderless" className="lg:w-80 shrink-0">
            <div className="flex flex-col items-center py-6 text-center">
              <div className="relative group">
                <Avatar
                  size={112}
                  src={avatar || undefined}
                  icon={<UserOutlined />}
                  style={{ background: '#b45309' }}
                />
                <Upload
                  showUploadList={false}
                  accept="image/*"
                  beforeUpload={handleUpload}
                  fileList={[] as UploadFile[]}
                >
                  <Button
                    size="small"
                    icon={<CameraOutlined />}
                    loading={uploading}
                    className="mt-3"
                  >
                    更换头像
                  </Button>
                </Upload>
              </div>
              <div className="font-display text-lg font-bold mt-4 text-[#1a1815]">
                {profile?.nickname || profile?.username || '—'}
              </div>
              <div className="text-xs text-[#8a7f6f] mt-1 tracking-wider">
                @{profile?.username || '—'}
              </div>
              <div className="w-8 h-[3px] bg-[#b45309] my-4" />
              <p className="m-0 text-sm text-[#6d6357] whitespace-pre-wrap leading-relaxed">
                {profile?.bio || '还没有写简介'}
              </p>
            </div>
          </Card>

          {/* 右侧：编辑表单 */}
          <Card variant="borderless" className="flex-1">
            <div className="max-w-xl pt-2">
              <Form form={form} layout="vertical">
                <Form.Item
                  name="nickname"
                  label="昵称"
                  rules={[
                    { max: 50, message: '昵称不能超过 50 个字符' },
                  ]}
                >
                  <Input placeholder="展示在站点的昵称" maxLength={50} showCount />
                </Form.Item>
                <Form.Item
                  name="bio"
                  label="个人简介"
                  rules={[{ max: 500, message: '简介不能超过 500 个字符' }]}
                >
                  <Input.TextArea
                    rows={5}
                    placeholder="写一段介绍自己的话…"
                    maxLength={500}
                    showCount
                  />
                </Form.Item>
                <Form.Item label="用户名">
                  <Input value={profile?.username} disabled />
                </Form.Item>
                <Form.Item>
                  <Button
                    type="primary"
                    icon={<SaveOutlined />}
                    loading={saving}
                    onClick={handleSave}
                  >
                    保存个人信息
                  </Button>
                </Form.Item>
              </Form>
            </div>
          </Card>
        </div>
      </Spin>
    </div>
  );
}
