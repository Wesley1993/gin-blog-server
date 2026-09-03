import { useEffect, useState } from 'react';
import {
  Avatar,
  Button,
  Card,
  Col,
  DatePicker,
  Divider,
  Form,
  Input,
  InputNumber,
  Row,
  Select,
  Space,
  Spin,
  Upload,
  App as AntApp,
} from 'antd';
import type { UploadFile } from 'antd';
import {
  CameraOutlined,
  DeleteOutlined,
  PlusOutlined,
  SaveOutlined,
  UserOutlined,
} from '@ant-design/icons';
import dayjs from 'dayjs';
import type { Dayjs } from 'dayjs';
import PageHeader from '../../components/PageHeader';
import { getSiteProfile, saveSiteProfile } from '../../api/siteProfile';
import { uploadToOss } from '../../api/upload';
import type { SiteProfile, SiteProfileProject, SiteProfileSkill } from '../../api/types';

/** 表单内部结构：项目起止时间用 Dayjs 承载 */
interface ProjectFormValue {
  name?: string;
  desc?: string;
  tech?: string[];
  start?: Dayjs | null;
  end?: Dayjs | null;
  link?: string;
}

interface ProfileFormValues {
  name?: string;
  title?: string;
  bio?: string;
  contacts?: { github?: string; email?: string; wechat?: string; qq?: string; address?: string };
  skills?: { name?: string; group?: string; level?: number }[];
  projects?: ProjectFormValue[];
}

/** 个人资料：编辑展示端「关于页 / 联系站长」使用的站长资料 */
export default function SiteProfilePage() {
  const { message } = AntApp.useApp();
  const [form] = Form.useForm<ProfileFormValues>();
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [avatar, setAvatar] = useState('');

  useEffect(() => {
    setLoading(true);
    getSiteProfile()
      .then((res) => {
        const data = res.data;
        setAvatar(data?.avatar ?? '');
        form.setFieldsValue({
          name: data?.name,
          title: data?.title,
          bio: data?.bio,
          contacts: {
            github: data?.contacts?.github,
            email: data?.contacts?.email,
            wechat: data?.contacts?.wechat,
            qq: data?.contacts?.qq,
            address: data?.contacts?.address,
          },
          skills: data?.skills ?? [],
          projects: (data?.projects ?? []).map((p) => ({
            name: p.name,
            desc: p.desc,
            tech: p.tech,
            start: p.start ? dayjs(p.start, 'YYYY-MM') : null,
            end: p.end ? dayjs(p.end, 'YYYY-MM') : null,
            link: p.link,
          })),
        });
      })
      .catch(() => undefined)
      .finally(() => setLoading(false));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  /** 头像选择后直接上传 OSS，成功后更新预览，随表单一起保存 */
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
    const payload: SiteProfile = {
      name: values.name ?? '',
      avatar,
      title: values.title ?? '',
      bio: values.bio ?? '',
      contacts: {
        github: values.contacts?.github ?? '',
        email: values.contacts?.email ?? '',
        wechat: values.contacts?.wechat ?? '',
        qq: values.contacts?.qq ?? '',
        address: values.contacts?.address ?? '',
      },
      skills: (values.skills ?? []).map((s): SiteProfileSkill => ({
        name: s?.name ?? '',
        level: typeof s?.level === 'number' ? s.level : 0,
        group: s?.group ?? '',
      })),
      projects: (values.projects ?? []).map((p): SiteProfileProject => ({
        name: p?.name ?? '',
        desc: p?.desc ?? '',
        tech: p?.tech ?? [],
        start: p?.start ? p.start.format('YYYY-MM') : '',
        end: p?.end ? p.end.format('YYYY-MM') : '',
        link: p?.link ?? '',
      })),
    };
    setSaving(true);
    try {
      await saveSiteProfile(payload);
      message.success('个人资料已保存');
    } catch {
      /* 拦截器已提示 */
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="page-container">
      <PageHeader title="个人资料" subtitle="站长资料，用于展示端「关于页」与「联系站长」" />

      <Spin spinning={loading}>
        <Form form={form} layout="vertical">
          <div className="flex flex-col gap-5">
            {/* 基础信息 */}
            <Card variant="borderless" title="基础信息">
              <div className="flex flex-col md:flex-row gap-6">
                <div className="flex flex-col items-center shrink-0 md:w-44 pt-1">
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
                    <Button size="small" icon={<CameraOutlined />} loading={uploading} className="mt-3">
                      更换头像
                    </Button>
                  </Upload>
                </div>
                <div className="flex-1 max-w-2xl">
                  <Row gutter={16}>
                    <Col xs={24} sm={12}>
                      <Form.Item name="name" label="姓名" rules={[{ max: 50, message: '姓名不能超过 50 个字符' }]}>
                        <Input placeholder="展示在关于页的姓名" maxLength={50} />
                      </Form.Item>
                    </Col>
                    <Col xs={24} sm={12}>
                      <Form.Item name="title" label="头衔" rules={[{ max: 100, message: '头衔不能超过 100 个字符' }]}>
                        <Input placeholder="一句话介绍，如：后端开发工程师" maxLength={100} />
                      </Form.Item>
                    </Col>
                    <Col span={24}>
                      <Form.Item name="bio" label="简历简介" extra="支持多段文本，展示于关于页">
                        <Input.TextArea rows={6} placeholder="写一段自我介绍…" />
                      </Form.Item>
                    </Col>
                  </Row>
                </div>
              </div>
            </Card>

            {/* 联系方式 */}
            <Card variant="borderless" title="联系方式">
              <Row gutter={16} className="max-w-3xl">
                <Col xs={24} sm={12}>
                  <Form.Item name={['contacts', 'github']} label="GitHub">
                    <Input placeholder="如：https://github.com/xxx" maxLength={100} allowClear />
                  </Form.Item>
                </Col>
                <Col xs={24} sm={12}>
                  <Form.Item name={['contacts', 'email']} label="邮箱">
                    <Input placeholder="如：hello@example.com" maxLength={100} allowClear />
                  </Form.Item>
                </Col>
                <Col xs={24} sm={12}>
                  <Form.Item name={['contacts', 'wechat']} label="微信">
                    <Input placeholder="微信号" maxLength={100} allowClear />
                  </Form.Item>
                </Col>
                <Col xs={24} sm={12}>
                  <Form.Item name={['contacts', 'qq']} label="QQ">
                    <Input placeholder="QQ 号" maxLength={100} allowClear />
                  </Form.Item>
                </Col>
                <Col xs={24} sm={12}>
                  <Form.Item name={['contacts', 'address']} label="地址">
                    <Input placeholder="请输入地址" maxLength={200} allowClear />
                  </Form.Item>
                </Col>
              </Row>
            </Card>

            {/* 技术栈 */}
            <Card
              variant="borderless"
              title="技术栈"
              extra={
                <span className="text-xs text-[#8a7f6f]">熟练度范围 0 - 100</span>
              }
            >
              <Form.List name="skills">
                {(fields, { add, remove }) => (
                  <>
                    {fields.map(({ key, name, ...rest }) => (
                      <Row key={key} gutter={12} align="middle" className="mb-2">
                        <Col xs={24} sm={8}>
                          <Form.Item
                            {...rest}
                            name={[name, 'name']}
                            className="mb-0"
                            rules={[{ required: true, message: '请输入技能名称' }]}
                          >
                            <Input placeholder="技能名称，如：Go" />
                          </Form.Item>
                        </Col>
                        <Col xs={10} sm={6}>
                          <Form.Item {...rest} name={[name, 'group']} className="mb-0">
                            <Input placeholder="分组，如：后端" />
                          </Form.Item>
                        </Col>
                        <Col xs={10} sm={6}>
                          <Form.Item {...rest} name={[name, 'level']} className="mb-0">
                            <InputNumber min={0} max={100} style={{ width: '100%' }} placeholder="熟练度" />
                          </Form.Item>
                        </Col>
                        <Col xs={4} sm={4}>
                          <Button
                            type="text"
                            danger
                            icon={<DeleteOutlined />}
                            onClick={() => remove(name)}
                          />
                        </Col>
                      </Row>
                    ))}
                    <Button
                      type="dashed"
                      block
                      icon={<PlusOutlined />}
                      onClick={() => add({ level: 60 })}
                    >
                      添加技能
                    </Button>
                  </>
                )}
              </Form.List>
            </Card>

            {/* 项目经历 */}
            <Card variant="borderless" title="项目经历">
              <Form.List name="projects">
                {(fields, { add, remove }) => (
                  <>
                    {fields.map(({ key, name, ...rest }) => (
                      <div
                        key={key}
                        className="border border-solid border-[#e7e0d2] rounded-md p-4 mb-4"
                      >
                        <div className="flex items-center justify-between mb-3">
                          <span className="text-xs tracking-[0.15em] text-[#8a7f6f] uppercase">
                            项目 #{name + 1}
                          </span>
                          <Button
                            type="text"
                            danger
                            size="small"
                            icon={<DeleteOutlined />}
                            onClick={() => remove(name)}
                          >
                            删除
                          </Button>
                        </div>
                        <Row gutter={16}>
                          <Col xs={24} sm={12}>
                            <Form.Item
                              {...rest}
                              name={[name, 'name']}
                              label="项目名称"
                              rules={[{ required: true, message: '请输入项目名称' }]}
                            >
                              <Input placeholder="如：博客系统" />
                            </Form.Item>
                          </Col>
                          <Col xs={24} sm={12}>
                            <Form.Item {...rest} name={[name, 'link']} label="项目链接">
                              <Input placeholder="https://…（可选）" allowClear />
                            </Form.Item>
                          </Col>
                          <Col span={24}>
                            <Form.Item {...rest} name={[name, 'desc']} label="项目描述">
                              <Input.TextArea rows={2} placeholder="项目做了什么、你的角色…" />
                            </Form.Item>
                          </Col>
                          <Col xs={24} sm={12}>
                            <Form.Item {...rest} name={[name, 'tech']} label="技术栈">
                              <Select
                                mode="tags"
                                placeholder="输入后回车添加，如：Go、Gin"
                                tokenSeparators={[',', '，']}
                                open={false}
                              />
                            </Form.Item>
                          </Col>
                          <Col xs={12} sm={6}>
                            <Form.Item {...rest} name={[name, 'start']} label="开始时间">
                              <DatePicker
                                picker="month"
                                format="YYYY-MM"
                                style={{ width: '100%' }}
                                placeholder="如：2026-01"
                              />
                            </Form.Item>
                          </Col>
                          <Col xs={12} sm={6}>
                            <Form.Item {...rest} name={[name, 'end']} label="结束时间" extra="留空表示至今">
                              <DatePicker
                                picker="month"
                                format="YYYY-MM"
                                style={{ width: '100%' }}
                                placeholder="至今"
                              />
                            </Form.Item>
                          </Col>
                        </Row>
                      </div>
                    ))}
                    <Button
                      type="dashed"
                      block
                      icon={<PlusOutlined />}
                      onClick={() => add({ tech: [] })}
                    >
                      添加项目
                    </Button>
                  </>
                )}
              </Form.List>
            </Card>
          </div>

          <Divider className="!my-6" />
          <div className="flex justify-end">
            <Space>
              <span className="text-xs text-[#8a7f6f]">保存后展示端将同步更新</span>
              <Button type="primary" size="large" icon={<SaveOutlined />} loading={saving} onClick={handleSave}>
                保存个人资料
              </Button>
            </Space>
          </div>
        </Form>
      </Spin>
    </div>
  );
}
