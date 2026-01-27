import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";

const fetchRules = async () => {
    const response = await fetch("http://localhost:10100/api/rules");
    if (!response.ok) throw new Error('Network response was not ok');
    return response.json();
};

const deleteRule = async (id) => {
    const response = await fetch(`http://localhost:10100/api/rules/${id}`, {
        method: 'DELETE',
    });
    if (!response.ok) throw new Error('Failed to delete device');
    return response.json();
};

export default function Rules() {
    const navigate = useNavigate();
    const queryClient = useQueryClient();

    const { data: rules, isLoading: loadingRules, error: errorRules } = useQuery({
        queryKey: ['rules'],
        queryFn: fetchRules,
    });

    const deleteRuleMutation = useMutation({
        mutationFn: deleteRule,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['rules'] })
        }
    })

    const handleRuleDelete = (id) => {
        deleteRuleMutation.mutate(id);
    }

    const handleRuleMore = (id) => {
        console.log(`more ${id}`);
    }

    const rulesToDisplay = rules?.map(rule => (
        <section key={rule.ruleID} className="rule-item">
            <span
                className="material-symbols-outlined rule-arrow rule-more"
                onClick={() => handleRuleMore(rule.ruleID)}
            >
                arrow_drop_down
            </span>
            <h1
                className="rule-more"
                onClick={() => handleRuleMore(rule.ruleID)}
            >
                Rule {rule.ruleID.split("-")[1]}
            </h1>
            <span
                className="material-symbols-outlined rule-delete"
                onClick={() => handleRuleDelete(rule.ruleID)}
            >
                delete
            </span>
        </section>
    ))

    return (
        <div className="rules-page">
            <section className="current-rules">
                {rules && rules.length > 0 ? rulesToDisplay : <h2>Found no rules. Click below to add some.</h2>}
            </section>
            <button
                className="add-rule-button add-button"
                onClick={() => navigate("/AddRules")}
            >
                Add a Rule
            </button>
        </div>
    )
}
