import { PageHeader } from '../ui';

export default function ComingSoon({ title, phase }: { title: string; phase: number }) {
  return (
    <>
      <PageHeader title={title} />
      <div className="card empty">
        <p>هذه الوحدة غير متاحة بعد، وسيتم بناؤها في <strong>المرحلة {phase}</strong>.</p>
        <p className="muted">لا تُعرض هنا أي بيانات تجريبية.</p>
      </div>
    </>
  );
}
