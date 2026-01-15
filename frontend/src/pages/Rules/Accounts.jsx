import { useState } from "react";

export default function Accounts(props) {
    const [accountsOpen, setAccountsOpen] = useState(false);
    const [selectedAccounts, setSelectedAccounts] = useState([]);

    const toggleAccount = (accountId) => {
        setSelectedAccounts(prev =>
            prev.includes(accountId)
                ? prev.filter(id => id !== accountId)
                : [...prev, accountId]
        );
    };

    return (
        <div
            className="select-accounts"
            onBlur={(e) => {
                if (!e.currentTarget.contains(e.relatedTarget)) {
                    setAccountsOpen(false);
                }
            }}
            tabIndex={0}
        >
            <div
                className="select-accounts-header"
                onClick={() => setAccountsOpen(prev => !prev)}
            >
                {selectedAccounts.length === 0 ? "Select Accounts..." : `(${selectedAccounts.length}) Account${selectedAccounts.length === 1 ? "" : "s"} Selected`} {accountsOpen ? '◀' : '▶'}
            </div>
            {accountsOpen && (
                <div
                    className="select-accounts-list"
                >
                    {props.data.map(account => (
                        <label key={account.id} className="select-account-item">
                            <input
                                type="checkbox"
                                checked={selectedAccounts.includes(account.id)}
                                onChange={() => toggleAccount(account.id)}
                            />
                            {account.title}
                        </label>
                    ))}
                </div>
            )}
        </div>
    )
}
