import "./PullsList.css"
import Pull from "./Pull";

function PullsList(props) {
    return (
        <div className="pulls-list">
            <Pull id={1} link={"www.google.com"}/>
        </div>
    );
}

export default PullsList;
