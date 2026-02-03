import { useState } from "react";

export default function Accounts(props) {
    const [accountsOpen, setAccountsOpen] = useState(false);

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
                {props.selectedAccounts.length === 0 ? "Select Accounts..." : `(${props.selectedAccounts.length}) Account${props.selectedAccounts.length === 1 ? "" : "s"} Selected`} {accountsOpen ? '◀' : '▶'}
            </div>
            {accountsOpen && (
                <div
                    className="select-accounts-list"
                >
                    {props.data && props.data.length > 0 ?
                        <>
                            <label className="select-account-item">
                                <input
                                    type="checkbox"
                                    checked={props.selectedAccounts.length === props.data.length}
                                    onChange={() => props.addAccount("toggleAll")}
                                />
                                Select All
                            </label>
                            {props.data.map(account => (
                                <label key={account.id} className="select-account-item">
                                    <input
                                        type="checkbox"
                                        checked={props.selectedAccounts.includes(account.id)}
                                        onChange={() => props.addAccount(account.id)}
                                    />
                                    {account.title}
                                </label>
                            ))}
                        </>
                        : <section className="rules-empty-array">No Accounts Saved</section>
                    }
                </div>
            )}
        </div>
    )
}
