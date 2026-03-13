import { ApprovalModal } from '../../components/ApprovalModal';

export default function ApprovalsPage() {
    return (
        <div>
            <h1 className="text-2xl font-bold mb-6">Approval Queue</h1>
            <p className="text-cf-muted mb-4">Pending approval gates will appear here.</p>
            <ApprovalModal />
        </div>
    );
}
