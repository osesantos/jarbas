import './MainView.css';
import RefreshButton from './RefreshButton';
import PullsList from "./PullsList";

function RenderList(){
    return(
        <div className="pulls-list-container">
            <PullsList />
        </div>
    );
}

function MainView() {
  return (
    <div className="main-view">
      <header className="main-view-header">
        <div className="main-view-container">
            <RefreshButton className="refresh-button"/>
            {RenderList()}
        </div>
      </header>
    </div>
  );
}

export default MainView;