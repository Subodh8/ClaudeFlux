import { CostHeatmap } from '../../components/CostHeatmap';

export default function CostsPage() {
    return (
        <div>
            <h1 className="text-2xl font-bold mb-6">Cost Breakdown</h1>
            <CostHeatmap />
        </div>
    );
}
