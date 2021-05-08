import './MainView.css';
import RefreshButton from './RefreshButton';

function MainView() {
  return (
    <div className="main-view">
      <header className="main-view-header">
        <div className="main-view-container">
            <RefreshButton/>
        </div>
      </header>
    </div>
  );
}

export default MainView;