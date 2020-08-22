import React from 'react'
import { Pane, Button, Heading, Badge, Spinner } from 'evergreen-ui'
import '../../components/container/style/card.css'

import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import { genericImage, runImageToContainer, toggleImageDeleteModal } from '../../store/actions/image.action'
import CreatedAt from '../createdAt'

class ImageCard extends React.PureComponent {

  bytesToString(value) {
    var metrics = ["B", "KB", "MB", "GB"]
    var str = value + metrics[0]
    for (var i = 0; i < metrics.length; i++) {
      value = Math.floor(value / 1024)
      if (value == 0) break      
      str = value + metrics[i+1]
    }
    return str
  }

  render () {
    const { image, activeIndex, genericImage, index, toggleImageDeleteModal, runImageToContainer } = this.props
    const active = activeIndex == index
   return <Pane 
            display="flex" 
            flexDirection="column" 
            flexGrow={1}
            padding={12}
            borderRadius={6}
            border="default"
            className={active ? "element-card card-active" : "element-card"}
            onMouseEnter={() => genericImage({
            activeIndex: index
            })}>
            <Pane display="flex">
               <Pane display="flex" justifyContent="center" alignItems="center">
                  <Heading size={400}>{`${image.tags != ''? image.tags : 'No Repository'}`}</Heading>
               </Pane>
               <Badge backgroundColor="#e7e9ef" fontWeight="bold" borderRadius={16} paddingLeft={10} fontSize={11} paddingRight={10} marginLeft={10} marginTop={3}>{image.id.replace("sha256:","").slice(0, 16)}</Badge>
               <Badge backgroundColor="#d4eee3" fontWeight="bold" borderRadius={16} paddingLeft={10} fontSize={11} paddingRight={10} marginLeft={10} marginTop={3}>{this.bytesToString(image.size)}</Badge>
               <CreatedAt time={image.created} />
            </Pane>
            { active && 
               <Pane display="flex" marginTop={12}>
                  <Button marginRight={5} 
                        height={22} 
                        iconBefore="application"
                        onClick={()=>{
                          runImageToContainer(image)
                        }}
                        isLoading={image.stateToggling}
                       >Run</Button>
                  <Button marginRight={5} 
                        height={22} 
                        iconBefore="trash" 
                        onClick={()=>{
                          toggleImageDeleteModal(image)
                        }}
                        >
                        Delete
                  </Button>
               </Pane>
            }
         </Pane>
  }
}

const mapStateToProps = state => {
  return {
    activeIndex: state.image.activeIndex,
  }
}

const mapDispatchToProps = dispatch => bindActionCreators(
  {
    genericImage,
    toggleImageDeleteModal,
    runImageToContainer
  },
  dispatch
)

export default connect(mapStateToProps, mapDispatchToProps)( ImageCard )