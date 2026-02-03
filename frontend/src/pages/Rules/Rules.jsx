import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
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

const fetchAllBulbs = async () => {
    const response = await fetch("http://localhost:10100/api/bulbs");
    if (!response.ok) throw new Error('Network response was not ok');
    return response.json();
}

export default function Rules() {
    const [rulesToShow, setRulesToShow] = useState([]);

    const navigate = useNavigate();
    const queryClient = useQueryClient();

    const { data: rules, isLoading: loadingRules, error: errorRules } = useQuery({
        queryKey: ['rules'],
        queryFn: fetchRules,
    });

    const { data: allBulbs, isLoading: loadingAllBulbs, error: errorBulbs } = useQuery({
        queryKey: ['bulbs'],
        queryFn: fetchAllBulbs,
    });

    const deleteRuleMutation = useMutation({
        mutationFn: deleteRule,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['rules'] })
        }
    });

    const handleRuleDelete = (id) => {
        deleteRuleMutation.mutate(id);
    };

    const handleRuleMore = (id) => {
        if (rulesToShow.includes(id)) {
            setRulesToShow(prev => prev.filter(ruleId => ruleId !== id));
            return;
        }
        setRulesToShow(prev => [...prev, id]);
        return;
    };

    const paramsToArray = (params) => {
        const paramsArr = [];
        for (const [cmd, value] of Object.entries(params)) {
            switch (cmd) {
                case "state":
                    paramsArr.push(`Turn ${value ? "On" : "Off"}`);
                    break;
                case "dimming":
                    paramsArr.push(`Set Dim to ${value}`);
                    break;
                case "temp":
                    paramsArr.push(`Set Temperature to ${value}`);
                    break;
                case "r":
                    paramsArr.push(`Set Red to ${value}`);
                    break;
                case "g":
                    paramsArr.push(`Set Green to ${value}`);
                    break;
                case "b":
                    paramsArr.push(`Set Blue to ${value}`);
                    break;
            }
        }
        return paramsArr;
    }

    console.log(rules);

    const rulesToDisplay = rules?.map(rule => (
        <div key={rule.ruleID} className="rule-item">
            <section className="rule-header">
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
            <section
                className={rulesToShow.includes(rule.ruleID) ? "rule-description-show" : "rule-description-hide"}
            >
                <h2>When User{!rule.condition.account || rule.condition.account.length <= 1 ? "" : "s"}:</h2>
                <section className="rule-accounts">
                    {!rule.condition.account ? <h3>Any</h3> : rule.condition.account.map(acc => (
                        <h3
                            key={acc.id}
                            className="rule-account-item"
                        >
                            {acc.title}
                        </h3>
                    ))}
                </section>
                <h2>Do{!rule.condition.account || rule.condition.account.length <= 1 ? "es" : ""} Event{rule.condition.event.length === 1 ? "" : "s"}:</h2>
                <section className="rule-events">
                    {rule.condition.event.map((evnt, idx) => (
                        <h2
                            key={idx}
                            className="rule-event-item"
                        >
                            {evnt.split(".")[1]}
                        </h2>
                    ))}
                </section>
                <h2>On Device{!rule.condition.device || rule.condition.device.length <= 1 ? "" : "s"}:</h2>
                <section className="rule-devices">
                    {!rule.condition.account ? <h3>Any</h3> : rule.condition.device.map(dev => (
                        <h3
                            key={dev.id}
                            className="rule-device-item"
                        >
                            {dev.name}
                        </h3>
                    ))}
                </section>
                <h2>Perform Command{paramsToArray(rule.action[0].command.params).length === 1 ? "" : "s"}:</h2>
                <section className="rule-commands">
                    {paramsToArray(rule.action[0].command.params).map((cmd, idx) => (
                        <h3
                            key={idx}
                            className="rule-command-item"
                        >
                            {cmd}
                        </h3>
                    ))}
                </section>
                <h2>To Bulb{rule.action[0].bulbsMac.length === 1 ? "" : "s"}:</h2>
                <section className="rule-bulbsMacs">
                    {rule.action[0].bulbsMac.map(mac => (
                        <h3
                            key={mac}
                            className="rule-bulbsMac-item"
                        >
                            {allBulbs ? allBulbs.find(bulb => bulb.mac === mac).name : mac}
                        </h3>
                    ))}
                </section>
            </section>
        </div>
    ));

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
