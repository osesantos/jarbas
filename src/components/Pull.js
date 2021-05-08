import "./Pull.css"

function Pull(props) {
    return (
        <div className={"pull-item id-" + props.id}>
            <a href={props.link}>{props.link}</a>
        </div>
    );
}

export default Pull;
