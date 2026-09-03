import { useEffect } from 'react';
import { Tag } from 'antd';
import {
  EnvironmentOutlined,
  GithubOutlined,
  LinkOutlined,
  MailOutlined,
} from '@ant-design/icons';
import { useSiteStore } from '../store/site';
import { useProfileStore } from '../store/profile';

export default function About() {
  const siteDesc = useSiteStore((s) => s.siteDesc);
  const profile = useProfileStore((s) => s.profile);
  const fetchProfile = useProfileStore((s) => s.fetchProfile);

  // 优先取接口数据，失败/为空时 store 内部回退静态资料
  useEffect(() => {
    fetchProfile();
  }, [fetchProfile]);

  // 技能按组归类
  const groups = Array.from(new Set(profile.skills.map((s) => s.group)));

  return (
    <div className="container about-page">
      <div className="section-head">
        <span className="zh">关于</span>
        <span className="rule" />
        <span className="meta">ABOUT THE AUTHOR</span>
      </div>

      <div className="about-grid">
        {/* 左：名片 */}
        <div className="profile-card">
          <div className="avatar">{profile.name.slice(0, 1)}</div>
          <h2>{profile.name}</h2>
          <div className="title">{profile.title}</div>
          <div className="profile-contact">
            <span>
              <EnvironmentOutlined /> {profile.location}
            </span>
            {profile.email && (
              <a href={`mailto:${profile.email}`}>
                <MailOutlined /> {profile.email}
              </a>
            )}
            {profile.github && (
              <a href={profile.github} target="_blank" rel="noreferrer">
                <GithubOutlined /> GitHub
              </a>
            )}
          </div>
        </div>

        {/* 右：简介 / 技术栈 / 项目经历 */}
        <div>
          <section className="about-section">
            <div className="side-title">
              <span className="zh">自序</span>
              <span className="en">Biography</span>
            </div>
            <p className="bio-text">
              {profile.bio}
              {siteDesc && (
                <>
                  <br />
                  <br />
                  关于本站：{siteDesc}
                </>
              )}
            </p>
          </section>

          <section className="about-section">
            <div className="side-title">
              <span className="zh">技术栈</span>
              <span className="en">Skills</span>
            </div>
            {groups.map((g) => (
              <div key={g}>
                <div className="skill-group-label">/ {g}</div>
                {profile.skills
                  .filter((s) => s.group === g)
                  .map((s) => (
                    <div className="skill-row" key={s.name}>
                      <span className="name">{s.name}</span>
                      <span className="bar">
                        <i style={{ width: `${s.level}%` }} />
                      </span>
                      <span className="num">{s.level}%</span>
                    </div>
                  ))}
              </div>
            ))}
          </section>

          <section className="about-section">
            <div className="side-title">
              <span className="zh">项目经历</span>
              <span className="en">Projects</span>
            </div>
            {profile.projects.map((p) => (
              <div className="project-card" key={p.name}>
                <div className="head">
                  <h3>
                    {p.link ? (
                      <a href={p.link} target="_blank" rel="noreferrer">
                        {p.name} <LinkOutlined />
                      </a>
                    ) : (
                      p.name
                    )}
                  </h3>
                  <span className="period">{p.period}</span>
                </div>
                <p className="desc">{p.desc}</p>
                <div className="tech">
                  {p.tech.map((t) => (
                    <Tag key={t} bordered={false}>
                      {t}
                    </Tag>
                  ))}
                </div>
              </div>
            ))}
          </section>
        </div>
      </div>
    </div>
  );
}
